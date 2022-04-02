package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/NorthfieldIT/yaml2confluence/internal/confluence"
)

func NewInstanceWizard() confluence.InstanceConfig {
	instance := confluence.InstanceConfig{}

	hostPrompt := &survey.Input{Message: "Confluence Host"}
	survey.AskOne(hostPrompt, &instance.Host, survey.WithValidator(survey.Required))
	var questions = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Instance Name",
				Default: strings.Split(instance.Host, ".")[0],
			},
			Validate: survey.Required,
		},
		{
			Name: "user",
			Prompt: &survey.Input{
				Message: "Username",
			},
			Validate: survey.Required,
		},
		{
			Name: "api_token",
			Prompt: &survey.Password{
				Message: "API Token",
			},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(questions, &instance)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return instance
}
