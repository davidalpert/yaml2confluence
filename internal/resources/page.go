package resources

import (
	"sort"

	"github.com/NorthfieldIT/yaml2confluence/internal/constants"
)

type ChangeType int

const (
	CREATE ChangeType = iota
	UPDATE
	DELETE
	NOOP
)

type Page struct {
	Key      string
	Resource *YamlResource
	Content  PageContent
	Remote   *RemoteResource
	Parent   *Page
	// childrenByTitle map[string]*Page
	Children []*Page
}

type PageContent struct {
	Markup string
	Sha256 string
}

func NewPage(key string, yr *YamlResource) *Page {
	return &Page{Key: key, Resource: yr}
}

func (p *Page) IsRoot() bool {
	return p.Parent == nil
}

func (p *Page) GetParent() *Page {
	return p.Parent
}

func (p *Page) GetChildren() []*Page {
	return p.Children
}

func (p *Page) GetChildByTitle(title string) *Page {
	for _, child := range p.Children {
		if child.GetTitle() == title {
			return child
		}
	}
	return nil
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
	parent.Children = append(parent.Children, p)

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

func (p *Page) LabelsDiffer() bool {
	// if we don't have labels to compare on both sides, return false
	if p.Resource == nil || p.Remote == nil {
		return false
	}

	l := append([]string{constants.GENERATED_BY_LABEL}, p.Resource.GetLabels()...)
	r := append([]string{}, p.Remote.Labels...)

	sort.Strings(l)
	sort.Strings(r)

	if len(l) != len(r) {
		return true
	}

	for i := range l {
		if l[i] != r[i] {
			return true
		}
	}

	return false
}
func (p *Page) Sha256Differs() bool {
	// if we don't have sha256s to compare on both sides, return false
	if p.Resource == nil || p.Remote == nil {
		return false
	}

	if p.Content.Sha256 != p.Remote.Sha256.Value {
		return true
	}

	return false
}
func (p *Page) GetChangeType() ChangeType {
	if p.Sha256Differs() || p.LabelsDiffer() {
		return UPDATE
	}
	if p.Resource != nil && p.Remote == nil {
		return CREATE
	}
	if p.Resource == nil && p.Remote != nil {
		return DELETE
	}

	return NOOP
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
	return p.Resource.Title
}
func (p *Page) GetAncestorId() string {
	return p.GetParent().GetRemoteId()
}
func (p *Page) GetContent() string {
	return p.Content.Markup
}
func (p *Page) GetLabels() []string {
	return p.Resource.GetLabels()
}
func (p *Page) GetIncrementedVersion() int {
	return p.GetRemoteVersion() + 1
}
func (p *Page) IsUpdate() bool {
	return p.GetChangeType() == UPDATE
}
