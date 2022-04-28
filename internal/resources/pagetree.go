package resources

import (
	"sort"
)

type PageTree struct {
	rootPage *Page
	pages    map[string]*Page
	deletes  [][]PageUpdate
}

type PageUpdate struct {
	Operation ChangeType
	Page      *Page
}

func NewPageTree(yr []*YamlResource, anchor string) *PageTree {
	pageTree := &PageTree{}

	pageTree.rootPage = NewPage("/", nil)
	if anchor != "" {
		pageTree.rootPage.Remote = &RemoteResource{Id: anchor}
	}
	pageTree.pages = map[string]*Page{
		"/": pageTree.rootPage,
	}

	for _, r := range yr {
		pageTree.AddPage(r)
	}

	return pageTree
}

type orphanPage struct {
	depth  int
	remote *RemoteResource
}

func (pt *PageTree) AddRemotes(remotes []*RemoteResource) {
	orphans := []orphanPage{}

	for _, remote := range remotes {
		titlePath := remote.GetTitlePath()
		page := pt.GetPageFromTitlePath(titlePath)
		if page != nil {
			page.Remote = remote
		} else {
			orphans = append(orphans, orphanPage{
				depth:  len(titlePath),
				remote: remote,
			})
		}
	}

	pt.setDeletes(orphans)
}

func (pt *PageTree) setDeletes(orphans []orphanPage) {
	if len(orphans) == 0 {
		return
	}

	sort.SliceStable(orphans, func(i, j int) bool {
		return orphans[i].depth > orphans[j].depth
	})

	deletes := [][]PageUpdate{}

	deletes = append(deletes, []PageUpdate{createDeletePageUpdate(orphans[0].remote)})
	currentDepth := orphans[0].depth

	for _, orphan := range orphans[1:] {
		if orphan.depth < currentDepth {
			deletes = append(deletes, []PageUpdate{createDeletePageUpdate(orphan.remote)})
		} else {
			deletes[len(deletes)-1] = append(deletes[len(deletes)-1], createDeletePageUpdate(orphan.remote))
		}
	}

	pt.deletes = deletes
}

func (pt *PageTree) AddPage(yr *YamlResource) {
	page := NewPage(yr.Path, yr)
	pt.pages[yr.GetParentPath()].AppendChild(page)
	pt.pages[page.Key] = page
}

func (pt *PageTree) GetPage(key string) *Page {
	return pt.pages[key]
}

func (pt *PageTree) GetPageFromTitlePath(titles []string) *Page {
	page := pt.rootPage
	for _, title := range titles {
		page = page.GetChildByTitle(title)
		if page == nil {
			return nil
		}
	}

	return page
}

func (pt *PageTree) GetPages() []*Page {
	pages := []*Page{}
	for _, p := range pt.pages {
		if p.IsRoot() {
			continue
		}
		pages = append(pages, p)
	}

	return pages
}

func (pt *PageTree) GetLevels() [][]string {
	levels := [][]string{}
	level := pt.rootPage.GetChildren()

	for len(level) > 0 {
		levels = append(levels, getKeys(level))

		children := []*Page{}
		for _, leaf := range level {
			children = append(children, leaf.GetChildren()...)
		}

		level = children
	}

	return levels
}

func (pt *PageTree) GetChanges() [][]PageUpdate {
	changes := [][]PageUpdate{}
	updates := []PageUpdate{}
	skips := []PageUpdate{}

	level := pt.rootPage.GetChildren()

	for len(level) > 0 {
		creates := []PageUpdate{}

		children := []*Page{}
		for _, page := range level {
			pu := createPageUpdate(page)
			switch pu.Operation {
			case CREATE:
				creates = append(creates, pu)
			case UPDATE:
				updates = append(updates, pu)
			case NOOP:
				skips = append(skips, pu)
			}

			children = append(children, page.GetChildren()...)
		}

		if len(creates) > 0 {
			changes = append(changes, creates)
		}

		level = children
	}

	// all updates can be applied in the first grouping
	changes = append(mergePageUpdates(pt.deletes, [][]PageUpdate{updates}), changes...)
	changes = append(changes, skips)

	return changes
}

func mergePageUpdates(p1, p2 [][]PageUpdate) [][]PageUpdate {
	m := [][]PageUpdate{}

	p1Length := len(p1)
	p2Length := len(p2)

	if p1Length == 0 && p2Length == 0 {
		return m
	} else if p1Length == 0 {
		return p2
	} else if p2Length == 0 {
		return p1
	}

	maxLength := p1Length
	if p2Length > maxLength {
		maxLength = p2Length
	}

	for i := 0; i < maxLength; i++ {
		v1 := []PageUpdate{}
		v2 := []PageUpdate{}

		if i < p1Length {
			v1 = p1[i]
		}
		if i < p2Length {
			v2 = p2[i]
		}

		m = append(m, append(v1, v2...))
	}

	return m
}

func createPageUpdate(p *Page) PageUpdate {
	return PageUpdate{
		Operation: p.GetChangeType(),
		Page:      p,
	}
}
func createDeletePageUpdate(remote *RemoteResource) PageUpdate {
	return PageUpdate{
		Operation: DELETE,
		Page:      &Page{Key: remote.Id, Remote: remote},
	}
}

func getKeys(nodes []*Page) []string {
	vals := []string{}
	for _, n := range nodes {
		vals = append(vals, n.Key)
	}

	return vals
}
