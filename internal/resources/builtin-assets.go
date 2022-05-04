package resources

var templates map[string]string = map[string]string{
	"wiki": "{{{markup}}}",
}

var hooks map[string]string = map[string]string{
	"required-fields": `target: .*
priority: -1
defaults:
   kind: wiki
   labels: []
yq: 
   # Defaults 'title' to the relative path (stored in head_comment)
   - '{} as $d|$d.title = (head_comment | capture("(.*)") .[])|. *n $d'
   # Defaults 'editorVersion' to instance setting (stored in foot_comment)
   - '{} as $d|$d.editorVersion = (foot_comment | capture("(.*)") .[])|. *n $d'
   # Remove foot_comment (editorVersion)
   - '. foot_comment=""'`,
}

type builtinAsset struct {
	name string
	data string
}

func (a builtinAsset) GetName() string {
	return a.name
}
func (a builtinAsset) GetPath() string {
	return ""
}
func (a builtinAsset) IsBuiltin() bool {
	return true
}
func (a builtinAsset) ReadString() string {
	return a.data
}

func (a builtinAsset) ReadBytes() []byte {
	return []byte(a.data)
}

func GetBuiltinTemplates() []IAsset {
	return toAssets(templates)
}

func GetBuiltinHooks() []IAsset {
	return toAssets(hooks)
}

func toAssets(assetMap map[string]string) []IAsset {
	assets := []IAsset{}
	for name, data := range assetMap {
		assets = append(assets, builtinAsset{
			name: name,
			data: data,
		})
	}

	return assets
}
