// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/manifoldco/promptui"
	assumeAccount "github.com/s3than/assume/account"
	"github.com/spf13/cobra"
)

type accountType struct {
	Name        string
	Description string
	Definition  string
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		addAccount()
	},
}

func addAccount() {

	accountTypes := []accountType{
		{
			Name:        "Base Account",
			Description: "AWS access keys are required.",
			Definition: `The combination of an access key ID (like AKIAIOSFODNN7EXAMPLE)
	and a secret access key (like wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY) are
	required for a base account.`,
		},
		{
			Name:        "Cross Account",
			Description: "An account with cross-account access, a Base Account is required.",
			Definition: `The process of permitting limited, controlled use of resources in
	one AWS account by a user in another AWS account.`,
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\u21D2 {{ .Name | cyan }} ({{ .Description | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Description | red }})",
		Selected: "\u21D2 {{ .Name | red | cyan }}",
		Details: `
--------- Account Details ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Definition:" | faint }}	{{ .Definition }}`,
	}

	searcher := func(input string, index int) bool {
		account := accountTypes[index]
		name := strings.Replace(strings.ToLower(account.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Add AWS Account",
		Items:     accountTypes,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if i == 0 {
		promptBaseAccount()
	} else if i == 1 {
		promptCrossAccount()
	}
}

func promptBaseAccount() {

	accounts, err := assumeAccount.FindAllbyType("base")

	if err != nil {
		panic("Unable to unmarshal hosts")
	}

	fmt.Printf("%+v", accounts)

}

func promptCrossAccount() {

	validateRoleName := func(input string) error {

		err := validation.Validate(input,
			validation.Required, // not empty
			validation.Match(regexp.MustCompile("^[A-Za-z0-9_=,.@+-]*$")),
			is.ASCII,
		)

		if err != nil {
			return err
		}
		return nil
	}

	validateAWSAccount := func(input string) error {

		err := validation.Validate(input,
			validation.Required, // not empty
			validation.Length(12, 12),
			is.UTFNumeric,
		)

		if err != nil {
			return err
		}
		return nil
	}

	promptProfile := promptui.Prompt{
		Label:    "Profile Name",
		Validate: validateRoleName,
	}

	profileName, err := promptProfile.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	promptRole := promptui.Prompt{
		Label:    "AWS Cross Account Role Name",
		Validate: validateRoleName,
	}

	roleName, err := promptRole.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	promptAccount := promptui.Prompt{
		Label:    "AWS Account Number",
		Validate: validateAWSAccount,
	}

	awsAccount, err := promptAccount.Run()

	accounts, err := assumeAccount.FindAllbyType("base")

	if err != nil {
		panic("Unable to unmarshal hosts")
	}

	var accountTypes []accountType

	for _, a := range accounts {
		accountTypes = append(accountTypes,
			accountType{
				Name: a.ProfileName,
			},
		)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\u21D2 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\u21D2 {{ .Name | red | cyan }}",
		Details: `
	--------- Base Account ----------
	{{ "Name:" | faint }}	{{ .Name }}`,
	}

	searcher := func(input string, index int) bool {
		account := accountTypes[index]
		name := strings.Replace(strings.ToLower(account.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Base Account Profile to link",
		Items:     accountTypes,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	newAccount := assumeAccount.Account{
		ProfileName:profileName,
		RoleArn: "arn:aws:iam::" + awsAccount + ":role/" + roleName,
		SourceProfile: accountTypes[i].Name,
	}

	assumeAccount.WriteAccountToConfig(newAccount)

}

func init() {
	rootCmd.AddCommand(addCmd)

	// fmt.Println("config")

	// subv := viper.GetStringMap("base_accounts.tcolbert")
	// fmt.Printf("%+v", subv)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
