package confluence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/NorthfieldIT/yaml2confluence/internal/constants"
)

var expectedStatusCode = map[string]int{
	"GET":    200,
	"POST":   200,
	"PUT":    200,
	"DELETE": 204,
}

type ConfluenceApiService struct {
	config   InstanceConfig
	spaceKey string
	client   http.Client
	authKey  string
	isCloud  bool
}
type NoOpResponse struct{}
type ConfluenceResponse interface {
	ConfluenceContentResponse | ConfluenceSearchResultsResponse | ConfluenceSpaceResponse | NoOpResponse
}

func NewConfluenceApiService(spaceKey string, config InstanceConfig) ConfluenceApiService {
	authKey := config.API_token
	isCloud := true
	if config.Type == "server" {
		authKey = config.Password
		isCloud = false
	}

	return ConfluenceApiService{
		config:   config,
		spaceKey: spaceKey,
		client: http.Client{
			Timeout: time.Second * 10,
		},
		authKey: authKey,
		isCloud: isCloud,
	}
}
func (api ConfluenceApiService) IsCloudInstance() bool {
	return api.isCloud
}
func (api ConfluenceApiService) IsServerInstance() bool {
	return !api.isCloud
}
func (api ConfluenceApiService) request(method string, URI string, body []byte) (*http.Response, error) {
	URL := api.config.Protocol + "://" + api.config.Host + filepath.Join(api.config.API_prefix, URI)
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(api.config.User, api.authKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expectedStatusCode[method] {
		body := ""
		if bodyText, err := ioutil.ReadAll(resp.Body); err == nil {
			body = string(bodyText)
		}
		return resp, errors.New(fmt.Sprintf("%s\n%s\n", resp.Status, body))
	}

	return resp, err
}

func unmarshallResponse[T ConfluenceResponse](resp *http.Response, err error) (T, error) {
	var result T

	if err != nil {
		return result, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(bodyText, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (api ConfluenceApiService) CreateSpaceIfNotExists() (bool, string, error) {
	// check is space exists already, if so, return
	content, err := unmarshallResponse[ConfluenceSpaceResponse](api.request("GET", fmt.Sprintf("/space/%s?expand=homepage", api.spaceKey), nil))
	if err == nil {
		return true, content.Homepage.Id, nil
	}

	payload := ConfluenceSpacePayload{
		Key:  api.spaceKey,
		Name: api.spaceKey,
	}

	postBody, _ := json.Marshal(payload)

	content, err = unmarshallResponse[ConfluenceSpaceResponse](api.request("POST", "/space/", postBody))
	if err != nil {
		return false, "", err
	}

	return false, content.Homepage.Id, nil

}

type UpsertPageContext interface {
	GetId() string
	GetTitle() string
	GetAncestorId() string
	GetContent() string
	GetLabels() []string
	GetIncrementedVersion() int
	IsUpdate() bool
}

func (api ConfluenceApiService) UpsertPage(page UpsertPageContext) (string, string, error) {
	method := "POST"
	uri := "/content"

	if page.IsUpdate() {
		method = "PUT"
		uri = uri + "/" + page.GetId()
	}

	labels := []Label{{Prefix: "global", Name: constants.GENERATED_BY_LABEL}}

	for _, l := range page.GetLabels() {
		labels = append(labels, Label{Prefix: "global", Name: l})
	}

	payload := ConfluenceContentPayload{
		Type:  "page",
		Title: page.GetTitle(),
		Space: Space{
			Key: api.spaceKey,
		},
		Version: Version{
			Number:    page.GetIncrementedVersion(),
			MinorEdit: true,
		},
		Body: Body{Storage{
			Value:          page.GetContent(),
			Representation: "wiki",
		}},
		Metadata: Metadata{
			Properties{Editor{
				Value: "V1",
			}},
			labels,
		},
	}

	if page.GetAncestorId() != "" {
		payload.Ancestors = append(payload.Ancestors, PageId{page.GetAncestorId()})
	}

	postBody, _ := json.Marshal(payload)

	content, err := unmarshallResponse[ConfluenceContentResponse](api.request(method, uri, postBody))
	if err != nil {
		fmt.Println(page.GetTitle())
		return "", "", err
	}

	return content.Id, content.Links.Base + content.Links.Webui, nil

}

func (api ConfluenceApiService) DeletePage(id string) error {
	_, err := api.request("DELETE", fmt.Sprintf("/content/%s", id), nil)

	// Confluence Server will mark a page as trashed on deletion
	// That page can then linger in the system for a few seconds, causing title collisions on page moves
	// Blindly attempt to permanently delete the trashed content, we don't care if this works or not
	if api.IsServerInstance() {
		api.request("DELETE", fmt.Sprintf("/content/%s?status=trashed", id), nil)
	}

	return err
}

type UpsertPropertyContext interface {
	GetId() string
	GetKey() string
	GetValue() string
	GetIncrementedVersion() int
	IsUpdate() bool
}

func (api ConfluenceApiService) UpsertProperty(property UpsertPropertyContext) error {
	method := "POST"
	if property.IsUpdate() {
		method = "PUT"
	}

	payload := ConfluenceContentPropertiesPayload{
		Value: property.GetValue(),
		Version: Version{
			Number:    property.GetIncrementedVersion(),
			MinorEdit: true,
		},
	}

	postBody, _ := json.Marshal(payload)

	_, err := api.request(method, fmt.Sprintf("/content/%s/property/%s", property.GetId(), property.GetKey()), postBody)

	return err
}

func (api ConfluenceApiService) SetLabels(contentId string, labels []string) error {
	payload := ConfluenceLabelPayload{}
	for _, label := range labels {
		payload = append(payload, Label{"global", label})
	}
	postBody, _ := json.Marshal(payload)
	_, err := api.request("POST", fmt.Sprintf("/content/%s/label", contentId), postBody)

	return err
}

func (api ConfluenceApiService) GetManagedContent() ([]ConfluencePageExpanded, string, error) {
	cql := url.PathEscape(fmt.Sprintf(`label="%s" AND space.key="%s"`, constants.GENERATED_BY_LABEL, api.spaceKey))
	URI := fmt.Sprintf("/content/search?cql=%s&expand=version,ancestors,metadata.properties.sha256,metadata.labels&limit=20", cql)

	sr, err := unmarshallResponse[ConfluenceSearchResultsResponse](api.request("GET", URI, nil))
	if err != nil {
		return nil, "", err
	}

	pages := sr.Results

	for sr.Links.Next != "" {
		sr, err := unmarshallResponse[ConfluenceSearchResultsResponse](api.request("GET", sr.Links.Context+sr.Links.Next, nil))
		if err != nil {
			return nil, "", err
		}

		pages = append(pages, sr.Results...)
	}

	return pages, sr.Links.Base, nil
}
