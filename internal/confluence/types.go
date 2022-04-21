package confluence

//---------------------
// PAYLOADS
//---------------------

// Space
type ConfluenceSpacePayload struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Content
type ConfluenceContentPayload struct {
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Space     Space    `json:"space"`
	Version   Version  `json:"version"`
	Ancestors []PageId `json:"ancestors"`
	Body      Body     `json:"body"`
}
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
type PageId struct {
	Id string `json:"id"`
}

// Label
type ConfluenceLabelPayload []Label
type Label struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
}

// Content Properties
type ConfluenceContentPropertiesPayload struct {
	Value   string  `json:"value"`
	Version Version `json:"version"`
}
type Version struct {
	Number    int  `json:"number"`
	MinorEdit bool `json:"minorEdit"`
}

//---------------------
// RESPONSES
//---------------------

// Content (Create and Update)
type ConfluenceContentResponse struct {
	Id    string
	Links struct {
		Webui string
		Base  string
	} `json:"_links"`
}

// Search Results
type ConfluenceSearchResultsResponse struct {
	Results []ConfluencePageExpanded
	Links   struct {
		Context string
		Next    string
		Base    string
	} `json:"_links"`
}
type ConfluencePageExpanded struct {
	ConfluencePage
	Version   Version
	Ancestors []ConfluencePage
	Metadata  struct {
		Properties struct {
			Sha256 struct {
				Id      string
				Value   string
				Version Version
			}
		}
	}
	Links struct {
		Webui string
	} `json:"_links"`
}
type ConfluencePage struct {
	Id    string
	Title string
}
