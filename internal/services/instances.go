package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"gopkg.in/yaml.v2"
)

type IInstancesSrv interface {
	New(string, bool)
	List()
}

type InstancesSrv struct{}

func NewInstancesService() InstancesSrv {
	return InstancesSrv{}
}

func (InstancesSrv) New(baseDir string, configOnly bool) {
	instance := cli.NewInstanceWizard()
	secret := utils.GetSecretAndGenerateIfMissing()

	var authPointer *string
	if instance.Type == "cloud" {
		authPointer = &instance.API_token
	} else {
		authPointer = &instance.Password
	}

	encryptedToken, err := utils.Encrypt(*authPointer, secret)
	if err != nil {
		panic(err)
	}
	*authPointer = "AES_ENC:" + encryptedToken

	configYaml, err := yaml.Marshal(&instance)
	if err != nil {
		panic(err)
	}

	if configOnly {
		fmt.Println("\n" + strings.TrimSpace(string(configYaml)))
		os.Exit(0)
	}

	utils.CreateInstanceDirectory(baseDir, instance.Name, string(configYaml))
}

func (InstancesSrv) List() {
	fmt.Println("-- none --")
}
