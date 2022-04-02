package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Space struct {
	Key string `json:"key"`
}
type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}
type Body struct {
	Storage Storage `json:"storage"`
}
type ConfluenceContentPayload struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Space Space  `json:"space"`
	Body  Body   `json:"body"`
}

func CreatePage(title string, content string, spaceKey string, config InstanceConfig) {
	payload := ConfluenceContentPayload{
		Type:  "page",
		Title: title,
		Space: Space{
			Key: spaceKey,
		},
		Body: Body{Storage{
			Value:          content,
			Representation: "wiki",
		}},
	}

	postBody, _ := json.Marshal(payload)

	client := &http.Client{}
	URL := "https://" + config.Host + "/wiki/rest/api/content"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(config.User, config.API_token)

	fmt.Println(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	fmt.Println(s)

}
