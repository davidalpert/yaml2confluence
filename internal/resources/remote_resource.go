package resources

type RemoteResource struct {
	Id        string
	Title     string
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

func (rr *RemoteResource) GetTitlePath() []string {
	titlePath := []string{}
	// append all of the ancestor titles, skipping the first one (the space page)
	for _, ancestor := range rr.Ancestors[1:] {
		titlePath = append(titlePath, ancestor.Title)
	}

	// add the title of the current page to the end
	titlePath = append(titlePath, rr.Title)

	return titlePath
}
