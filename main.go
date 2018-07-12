// Copyright © 2017 Tim Colbert admin@tcolbert.net
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

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

	flag "github.com/ogier/pflag"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/pquerna/otp/totp"
)

// AwsAccount struct to hold base account details
type AwsAccount struct {
	region          string
	profileName     string
	accessKeyID     string
	secretAccessKey string
}

// MfaToken token configuration
type MfaToken struct {
	serialNumber string
	tokenCode    string
}

// AssumeCredentials Credential details
type AssumeCredentials struct {
	profile  string
	secret   string
	account  AwsAccount
	roleArn  string
	session  client.ConfigProvider
	mfaToken MfaToken
}

var defaultConfig = "/.config/assume/config.ini"
var defaultCreds = "/.config/assume/config.creds"

var usr, err = user.Current()
var configFilePath = usr.HomeDir + "/.aws/config"
var credFilePath = usr.HomeDir + "/.aws/credentials"

func main() {

	var configCredsPath string
	var configPath string
	var accountRef string

	config := flag.StringP(
		"config",
		"s",
		"",
		"config file (default is $HOME/.config/assume/config.ini)",
	)
	configCreds := flag.StringP(
		"credentials",
		"c",
		"",
		"credentials file (default is $HOME/.config/assume/config.creds)",
	)

	// Account to assume, can use cli argument instead
	account := flag.StringP(
		"account",
		"a",
		"default",
		"AWS account reference",
	)

	saveProfile := flag.StringP(
		"profile",
		"p",
		"default",
		"AWS account profile to save as.",
	)

	flag.Parse()

	// Allow argument for account as well as -account

	if len(os.Args) > 1 {
		accountRef = os.Args[1]
		if strings.HasPrefix(accountRef, "-") {
			accountRef = *account
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *config == "" {
		configPath = usr.HomeDir + defaultConfig
	} else {
		configPath = *config
	}

	if *configCreds == "" {
		configCredsPath = usr.HomeDir + defaultCreds
	} else {
		configCredsPath = *config
	}

	os.OpenFile(configPath, os.O_CREATE, 0666)
	configFile, err := ini.Load(configPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.OpenFile(configCredsPath, os.O_CREATE, 0666)
	configFile.Append(configCredsPath)

	creds, err := getCredentials(configFile, accountRef, AssumeCredentials{})

	if creds.profile != "" {
		creds, err = getCredentials(configFile, creds.profile, creds)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if isConfigured(creds.account) {
		awsCreds := &sts.Credentials{}
		awsCreds, err = assumeRole(creds)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if awsCreds != nil {
			creds.account.profileName = *saveProfile
			account := creds.account

			err = writeFile(awsCreds, account)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			os.Exit(0)
		}
		os.Exit(1)
	}
}

func writeFile(awsCreds *sts.Credentials, account AwsAccount) error {
	os.OpenFile(credFilePath, os.O_CREATE, 0666)
	credFile, err := ini.Load(credFilePath)
	if err != nil {
		return err
	}
	credSect, err := credFile.NewSection(account.profileName)
	if err != nil {
		return err
	}
	credSect.NewKey("aws_access_key_id", *awsCreds.AccessKeyId)
	credSect.NewKey("aws_secret_access_key", *awsCreds.SecretAccessKey)
	credSect.NewKey("aws_session_token", *awsCreds.SessionToken)
	// Legacy support for boto2 apps
	credSect.NewKey("aws_security_token", *awsCreds.SessionToken)
	credSect.NewKey("region", account.region)
	credSect.NewKey("output", "json")

	credFile.SaveTo(credFilePath)

	// Config file details
	os.OpenFile(configFilePath, os.O_CREATE, 0666)
	configFile, err := ini.Load(configFilePath)
	if err != nil {
		return err
	}
	configSect, err := configFile.NewSection(account.profileName)
	if err != nil {
		return err
	}
	configSect.NewKey("region", account.region)
	configSect.NewKey("output", "json")
	configFile.SaveTo(configFilePath)

	return err
}

func assumeRole(creds AssumeCredentials) (*sts.Credentials, error) {

	mfaToken := creds.mfaToken

	var assumeRoleOutput *sts.AssumeRoleOutput        // nil
	var sessionTokenOutput *sts.GetSessionTokenOutput // nil
	var stsCredentials *sts.Credentials
	var err error

	switch {
	case mfaToken.serialNumber != "" && creds.profile != "":
		assumeRoleOutput, err = sts.New(creds.session).AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(creds.roleArn),
			RoleSessionName: aws.String(randStringBytesMaskImprSrc(6)),
			SerialNumber:    aws.String(mfaToken.serialNumber),
			TokenCode:       aws.String(mfaToken.tokenCode),
		})
		return assumeRoleOutput.Credentials, err
	case mfaToken.serialNumber != "" && creds.profile == "":
		sessionTokenOutput, err = sts.New(creds.session).GetSessionToken(&sts.GetSessionTokenInput{
			DurationSeconds:aws.Int64(43200),
			SerialNumber: aws.String(mfaToken.serialNumber),
			TokenCode:    aws.String(mfaToken.tokenCode),
		})
		return sessionTokenOutput.Credentials, err
	case mfaToken.serialNumber == "" && creds.profile != "":
		assumeRoleOutput, err = sts.New(creds.session).AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(creds.roleArn),
			RoleSessionName: aws.String(randStringBytesMaskImprSrc(6)),
		})
		return assumeRoleOutput.Credentials, err
	case mfaToken.serialNumber == "" && creds.profile == "":
		sessionTokenOutput, err = sts.New(creds.session).GetSessionToken(&sts.GetSessionTokenInput{
			DurationSeconds:aws.Int64(43200),
		})
		return sessionTokenOutput.Credentials, err
	}

	return stsCredentials, errors.New("no session set")
}

func getCredentials(config *ini.File, account string, creds AssumeCredentials) (AssumeCredentials, error) {
	section, err := config.GetSection(account)

	if section == nil {
		section, err = config.GetSection("profile " + account)
	}

	if section == nil {
		return creds, errors.New("no AWS config located")
	}

	profile, err := section.GetKey("source_profile")
	if profile != nil {
		creds.profile = profile.String()
	}

	secret, err := section.GetKey("secret")
	if secret != nil {
		creds.secret = secret.String()
	}

	awsAccount := AwsAccount{}

	region, err := section.GetKey("region")
	if region != nil {
		awsAccount.region = region.String()
	} else if region == nil {
		awsAccount.region = "ap-southeast-2"
	}

	awsAccessKeyID, err := section.GetKey("aws_access_key_id")
	if awsAccessKeyID != nil {
		awsAccount.accessKeyID = awsAccessKeyID.String()
	}

	awsSecretAccessKey, err := section.GetKey("aws_secret_access_key")
	if awsSecretAccessKey != nil {
		awsAccount.secretAccessKey = awsSecretAccessKey.String()
	}

	creds.account = awsAccount

	roleArn, err := section.GetKey("role_arn")

	if roleArn != nil {
		creds.roleArn = roleArn.String()
	}

	// Set AWS Session
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(awsAccount.region),
			Credentials: credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
				AccessKeyID:     awsAccount.accessKeyID,
				SecretAccessKey: awsAccount.secretAccessKey,
			}}),
		},
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	}))

	creds.session = awsSession

	devices, err := iam.New(creds.session).ListMFADevices(&iam.ListMFADevicesInput{})

	if len(devices.MFADevices) != 0 {
		token := ""
		awsToken := MfaToken{}
		if creds.secret == "" {
			token, err = stscreds.StdinTokenProvider()
		} else {
			token, err = totp.GenerateCode(creds.secret, time.Now())
		}
		awsToken.tokenCode = token
		awsToken.serialNumber = *devices.MFADevices[0].SerialNumber
		creds.mfaToken = awsToken
	}

	return creds, err
}

// isConfigured Return boolean if account details are set
func isConfigured(a AwsAccount) bool {

	if a.accessKeyID == "" || a.secretAccessKey == "" {
		return false
	}

	return true
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// randStringBytesMaskImprSrc generate random string
func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
