package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/NorthfieldIT/yaml2confluence/internal/confluence"
)

func NewInstanceWizard(baseDir string) (confluence.InstanceConfig, string) {
	instance := confluence.InstanceConfig{}

	confluenceTypeSelect := &survey.Select{
		Message: "Confluence Type",
		Options: []string{"cloud", "server"},
		Default: "cloud",
	}
	survey.AskOne(confluenceTypeSelect, &instance.Type, survey.WithValidator(survey.Required))
	hostPrompt := &survey.Input{Message: "Confluence Host"}
	survey.AskOne(hostPrompt, &instance.Host, survey.WithValidator(survey.Required))

	var err error
	if instance.Type == "cloud" {
		err = survey.Ask(getCloudQuestions(instance.Host), &instance)
	} else {
		err = survey.Ask(getServerQuestions(instance.Host), &instance)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var instanceDir string
	savePrompt := &survey.Select{
		Message: "Choose where to save instance configuration:",
		Options: []string{filepath.Join(baseDir, instance.Name), filepath.Join(baseDir)},
	}
	survey.AskOne(savePrompt, &instanceDir, survey.WithValidator(survey.Required))

	return instance, instanceDir
}

var user survey.Question = survey.Question{
	Name: "user",
	Prompt: &survey.Input{
		Message: "Username",
	},
	Validate: survey.Required,
}
var token survey.Question = survey.Question{
	Name: "api_token",
	Prompt: &survey.Password{
		Message: "API Token",
	},
	Validate: survey.Required,
}
var password survey.Question = survey.Question{
	Name: "password",
	Prompt: &survey.Password{
		Message: "Password",
	},
	Validate: survey.Required,
}

func getCloudQuestions(host string) []*survey.Question {
	return []*survey.Question{
		getNameQuestion(host),
		getApiPrefixQuestion("/wiki/rest/api"),
		&user,
		&token,
	}
}

func getServerQuestions(host string) []*survey.Question {
	return []*survey.Question{
		getNameQuestion(host),
		getApiPrefixQuestion("/rest/api"),
		&user,
		&password,
	}
}

func getNameQuestion(host string) *survey.Question {
	return &survey.Question{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Instance Name",
			Default: strings.Split(host, ".")[0],
		},
		Validate: survey.Required,
	}
}

func getApiPrefixQuestion(prefixDefault string) *survey.Question {
	return &survey.Question{
		Name: "api_prefix",
		Prompt: &survey.Input{
			Message: "API Base",
			Default: prefixDefault,
		},
		Validate: survey.Required,
	}
}
