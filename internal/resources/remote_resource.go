package resources

type RemoteResource struct {
	Id        string
	Title     string
	Labels    []string
	Link      string
	Version   int
	Ancestors []Ancestor
	Sha256    RemoteSha256
}

type Ancestor struct {
	Id    string
	Title string
}

type RemoteSha256 struct {
	Id      string
	Value   string
	Version int
}

func (rr *RemoteResource) GetTitlePath(anchorId string) []string {
	titlePath := []string{}
	startIndex := 1 // first page after space page

	if anchorId != "" {
		for i, ancestor := range rr.Ancestors[startIndex:] {
			if ancestor.Id == anchorId {
				startIndex = i + 2
			}
		}
	}
	// append all of the ancestor titles, skipping the space page, and possibly the hierachy upto the anchor
	for _, ancestor := range rr.Ancestors[startIndex:] {
		titlePath = append(titlePath, ancestor.Title)
	}

	// add the title of the current page to the end
	titlePath = append(titlePath, rr.Title)

	return titlePath
}
