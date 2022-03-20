package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/NorthfieldIT/yaml2confluence/internal/cli"
	"github.com/NorthfieldIT/yaml2confluence/internal/utils"
	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v2"
)

var instancesUsage = `
Usage:
	y2c instances new [<base_dir>] [--config-only]
	y2c instances list
	y2c instances decrypt <config>

Options:
	--config-only  	Skips directory creation, outputting the generated config file to the terminal
Sub-Commands:
	new  			A CLI wizard to configure a new Confluence instance
	list  			Lists all Confluence instances
	decrypt  		Displays a config yaml file with decrypted secrets
`

var instancesCmd = func(args docopt.Opts) {
	if args["list"].(bool) {
		fmt.Println("-- none --")
		return
	} else if args["new"].(bool) {
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

		if args["--config-only"].(bool) {

			fmt.Println("\n" + strings.TrimSpace(string(configYaml)))
			os.Exit(0)
		}
		// fmt.Println(encryptedToken)

		// decrypted, _ := utils.Decrypt(encryptedToken, secret)

		// fmt.Println(decrypted)

		dir := ""
		if args["<base_dir>"] != nil {
			dir = args["<base_dir>"].(string)
		}
		utils.CreateInstanceDirectory(dir, instance.Name, string(configYaml))

	}
}

func init() {
	cli.RegisterCommand("instances", instancesUsage, instancesCmd)
}
