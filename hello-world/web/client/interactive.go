package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{
			Message: "Insert your username",
			Help:    "It can be anything you like",
		},
		Validate: survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "password",
		Prompt: &survey.Password{
			Message: "Please type a strong password",
		},
		Validate: survey.Required,
	},
	{
		Name: "multifactor",
		Prompt: &survey.Confirm{
			Message: "Do you want to enable multifactor authentication?",
		},
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
}

func main2() {
	// the answers will be written to this struct
	answers := struct {
		Name          string                  // survey will match the question and field names
		FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
		Password      string
		Multifactor   bool
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)
	if answers.Multifactor {
		fmt.Println("multifactor is enabled")
	}

}
