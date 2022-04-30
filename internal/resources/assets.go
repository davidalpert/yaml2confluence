package resources

import (
	"os"
	"path/filepath"
	"strings"
)

type IAsset interface {
	GetName() string
	GetPath() string
	IsBuiltin() bool
	ReadString() string
	ReadBytes() []byte
}
type Asset struct {
	name string
	path string
}

func (a Asset) GetName() string {
	return a.name
}
func (a Asset) GetPath() string {
	return a.path
}
func (a Asset) IsBuiltin() bool {
	return false
}
func (a Asset) ReadString() string {
	return string(a.ReadBytes())
}

func (a Asset) ReadBytes() []byte {
	data, err := os.ReadFile(a.path)
	if err != nil {
		panic(err)
	}

	return data
}

func LoadAssets(dir string, exts []string) []IAsset {
	assets := []IAsset{}
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}

			if info.IsDir() {
				return nil
			}

			if !extensionIsOneOf(path, exts) {
				return nil
			}

			assets = append(assets, Asset{
				name: fileNameWithoutExtension(path),
				path: path,
			})

			return nil
		})
	if err != nil {
		panic(err)
	}

	return assets
}

func extensionIsOneOf(file string, exts []string) bool {
	for _, ext := range exts {
		if filepath.Ext(file) == ext {
			return true
		}
	}

	return false
}

func fileNameWithoutExtension(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
