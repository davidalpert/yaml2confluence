package confluence

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"gopkg.in/yaml.v2"
)

type InstanceConfig struct {
	Name      string
	Host      string
	User      string
	API_token string
}

type DirectoryProperties struct {
	ConfigPath   string
	SpaceDir     string
	TemplatesDir string
	SpaceKey     string
}

func LoadConfig(file string) InstanceConfig {
	ic := InstanceConfig{}

	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), &ic)
	if err != nil {
		panic(err)
	}

	secret, err := utils.GetSecret()
	if err != nil {
		fmt.Println("Could not find .secret")
		os.Exit(1)
	}
	decryptedToken, err := utils.Decrypt(strings.Split(ic.API_token, ":")[1], secret)
	if err != nil {
		fmt.Println("Could not decrypt API token")
		os.Exit(1)
	}

	ic.API_token = strings.TrimFunc(decryptedToken, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	return ic
}

func GetDirectoryProperties(path string) DirectoryProperties {
	dirTokens := strings.Split(utils.ResolveAbsolutePathFile(path), "spaces/")
	baseDir := dirTokens[0]
	spaceKey := strings.Split(dirTokens[1], "/")[0]

	props := DirectoryProperties{}
	props.ConfigPath = filepath.Join(baseDir, "config.yml")
	props.SpaceDir = filepath.Join(baseDir, "spaces", spaceKey)
	props.SpaceKey = spaceKey
	props.TemplatesDir = filepath.Join(baseDir, "templates")

	if _, err := os.Stat(props.ConfigPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Could not find config.yml")
		os.Exit(1)
	}

	if stat, err := os.Stat(props.SpaceDir); errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		fmt.Printf("Could not find '%s' space directory", props.SpaceKey)
		os.Exit(1)
	}

	if stat, err := os.Stat(props.TemplatesDir); errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		fmt.Println("Could not find templates directory")
		os.Exit(1)
	}

	return props
}
