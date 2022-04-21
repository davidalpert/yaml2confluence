package resources

type ChangeType int

const (
	CREATE ChangeType = iota
	UPDATE
	DELETE
	NOOP
)

type Page struct {
	Key             string
	Title           string
	Resource        *YamlResource
	Content         PageContent
	Remote          *RemoteResource
	Parent          *Page
	childrenByTitle map[string]*Page
}

type PageContent struct {
	Markup string
	Sha256 string
}

func NewPage(key string, yr *YamlResource) *Page {
	title := ""
	if yr != nil {
		title = yr.Title
	}
	return &Page{Key: key, Title: title, Resource: yr, childrenByTitle: map[string]*Page{}}
}

func (p *Page) IsRoot() bool {
	return p.Parent == nil
}

func (p *Page) GetParent() *Page {
	return p.Parent
}

func (p *Page) GetChildren() []*Page {
	children := []*Page{}
	for _, p := range p.childrenByTitle {
		children = append(children, p)
	}
	return children
}

func (p *Page) GetKeyArray() []string {
	keyArray := []string{p.Resource.Title}

	page := p
	for !page.GetParent().IsRoot() {
		page = page.GetParent()
		keyArray = append([]string{page.Resource.Title}, keyArray...)
	}

	return keyArray
}

func (parent *Page) AppendChild(p *Page) *Page {
	if parent == nil || p == nil || !p.IsRoot() {
		return nil
	}

	p.Parent = parent
	parent.childrenByTitle[p.Title] = p

	return p
}

func (p *Page) GetRemoteId() string {
	if p.Remote == nil {
		return ""
	}

	return p.Remote.Id
}

func (p *Page) GetRemoteVersion() int {
	if p.Remote == nil {
		return 0
	}

	return p.Remote.Version
}

func (p *Page) GetRemoteSha256Version() int {
	if p.Remote == nil {
		return 0
	}

	return p.Remote.Sha256.Version
}

func (p *Page) GetChangeType() ChangeType {
	if p.Resource != nil && p.Remote != nil {
		if p.Content.Sha256 != p.Remote.Sha256.Value {
			return UPDATE
		} else {
			return NOOP
		}
	}
	if p.Resource == nil && p.Remote != nil {
		return DELETE
	}

	return CREATE
}

func (p *Page) GetSha256Property() Property {
	return NewProperty(p.GetRemoteId(), "sha256", p.Content.Sha256, p.GetRemoteSha256Version())
}

// -------------------------
// UpsertContext functions
// -------------------------

func (p *Page) GetId() string {
	return p.GetRemoteId()
}
func (p *Page) GetTitle() string {
	return p.Title
}
func (p *Page) GetAncestorId() string {
	return p.GetParent().GetRemoteId()
}
func (p *Page) GetContent() string {
	return p.Content.Markup
}
func (p *Page) GetIncrementedVersion() int {
	return p.GetRemoteVersion() + 1
}
func (p *Page) IsUpdate() bool {
	return p.GetChangeType() == UPDATE
}
