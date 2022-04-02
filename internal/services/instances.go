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

	encryptedToken, err := utils.Encrypt(instance.API_token, secret)
	if err != nil {
		panic(err)
	}
	instance.API_token = "AES_ENC:" + encryptedToken

	configYaml, err := yaml.Marshal(&instance)
	if err != nil {
		panic(err)
	}

	if configOnly {

		fmt.Println("\n" + strings.TrimSpace(string(configYaml)))
		os.Exit(0)
	}
	// fmt.Println(encryptedToken)

	// decrypted, _ := utils.Decrypt(encryptedToken, secret)

	// fmt.Println(decrypted)

	utils.CreateInstanceDirectory(baseDir, instance.Name, string(configYaml))
}

func (InstancesSrv) List() {
	fmt.Println("-- none --")
}
