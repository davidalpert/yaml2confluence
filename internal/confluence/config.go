package confluence

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"gopkg.in/yaml.v3"
)

type InstanceConfig struct {
	Name       string
	Type       string
	Host       string
	API_prefix string
	User       string
	API_token  string `yaml:"api_token,omitempty"`
	Password   string `yaml:"password,omitempty"`
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
	var authPointer *string
	var errMsg string
	if ic.Type == "cloud" {
		authPointer = &ic.API_token
		errMsg = "Could not decrypt API token"
	} else {
		authPointer = &ic.Password
		errMsg = "Could not decrypt password"
	}
	decryptedToken, err := utils.Decrypt(strings.Split(*authPointer, ":")[1], secret)
	if err != nil {
		fmt.Println(errMsg)
		os.Exit(1)
	}

	*authPointer = strings.TrimFunc(decryptedToken, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	return ic
}
