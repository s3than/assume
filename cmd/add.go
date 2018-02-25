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
	"github.com/s3than/assume/account"
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

// type BaseAccount struct {
// 	ProfileName        string `mapstructure:"profile_name"`
// 	AwsAccessKeyID     string `mapstructure:"aws_access_key_id"`
// 	AwsSecretAccessKey string `mapstructure:"aws_secret_access_key"`
// 	Region             string
// 	Secret             string
// }

// type CrossAccount struct {
// 	ProfileName   string `mapstructure:"profile_name"`
// 	RoleArn       string `mapstructure:"role_arn"`
// 	SourceProfile string `mapstructure:"source_profile"`
// }

// type Accounts struct {
// 	BaseAccounts  []BaseAccount  `mapstructure:"base_accounts"`
// 	CrossAccounts []CrossAccount `mapstructure:"cross_accounts"`
// 	Hosts         []Host
// }

// type Host struct {
// 	Name string `mapstructure:"name"`
// 	Port int
// 	Key  string
// }

func promptBaseAccount() {

	account.FindAllbyType("base")

	// fmt.Printf("%+v", accounts)
	// config.getBaseAccounts()
	// var accounts Accounts

	// err := viper.Unmarshal(&accounts)

	// if err != nil {
	// 	panic("Unable to unmarshal hosts")
	// }

	// fmt.Printf("%+v", accounts)
	// var hosts []Host
	// err := viper.UnmarshalKey("hosts", &hosts)
	// if err != nil {
	// 	panic("Unable to unmarshal hosts")
	// }
	// for _, h := range hosts {
	// 	fmt.Printf("Name: %s, Port: %d, Key: %s\n", h.Name, h.Port, h.Key)
	// }
	// var baseAccounts []baseAccount
	// var B baseAccounts

	// baseAccountName := "tcolbert"

	// err := viper.UnmarshalKey("base_accounts", &baseAccounts)

	// if err != nil {
	// 	panic("Unable to unmarshal config")
	// }

	// fmt.Printf("%+v", baseAccounts)
	// for _, h := range B.baseAccounts {
	// 	fmt.Printf("profileName: %s, awsAccessKeyID: %s, awsSecretAccessKey: %s\n", h.profileName, h.awsAccessKeyID, h.awsSecretAccessKey)
	// }

	// subv := viper.Get("base_accounts." + baseAccountName)
	// Unmarshal(&B)

	// viper
	// viper.Set("base_accounts.test", subv)
	// viper.WriteConfig()
	// subv["test"] = "test"
	// fmt.Printf("%+v", B)

	// fmt.Printf("%+v", subv)
	// fmt.Printf("%+v", err)
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

	promptRole := promptui.Prompt{
		Label:    "AWS Role Name",
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

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("Your Aws Role is %q\n", roleName)
	fmt.Printf("Your Aws Account is %q\n", awsAccount)
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
