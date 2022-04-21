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
	Name      string
	Host      string
	User      string
	API_token string
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
