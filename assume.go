package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/pquerna/otp/totp"
)

type credentials struct {
	Expiration      time.Time `ini:"expiration,omitempty"`
	Duration        int64     `ini:"duration,omitempty"`
	Region          string    `ini:"region,omitempty"`
	AwsAccessKeyID  string    `ini:"aws_access_key_id,omitempty"`
	SecretAccessKey string    `ini:"aws_secret_access_key,omitempty"`
	MfaSecret       string    `ini:"secret,omitempty"`
	RoleArn         string    `ini:"role_arn,omitempty"`
	SourceProfile   string    `ini:"source_profile,omitempty"`
	SessionToken    string    `ini:"aws_session_token,omitempty"`
	SecurityToken   string    `ini:"aws_security_token,omitempty"`
	Output          string    `ini:"output,omitempty"`
	NamedProfile    string    `ini:"named_profile,omitempty"`
}

// Components encapsulate the individual pieces of an AWS ARN.
type components struct {
	ARN               string
	Partition         string
	Service           string
	Region            string
	AccountID         string
	ResourceType      string
	Resource          string
	ResourceDelimiter string
}

// MfaToken token configuration
type token struct {
	serialNumber string
	tokenCode    string
}

func assumeCommand(args arguments) {

	config, err := getCredentials(args)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cred, err := generateCredentials(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = writeFile(cred, config, args.saveProfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func getSection(cfg *ini.File, a string) (*ini.Section, error) {

	sect, err := cfg.GetSection(a)
	if err != nil {
		// Check for profile prefix in ini
		sect, err = cfg.GetSection("profile " + a)
	}

	return sect, err
}

func getCredentials(args arguments) (credentials, error) {
	var err error
	var a = args.account
	var c = credentials{
		NamedProfile: a,
	}
	cfg, err := ini.Load(configFile)
	if err != nil {
		return c, err
	}
	cfg.Append(credFile)

	sect, err := getSection(cfg, a)
	if err != nil {
		return c, err
	}
	err = sect.MapTo(&c)
	if c.SourceProfile != "" {
		sect, err := getSection(cfg, c.SourceProfile)
		if err != nil {
			return c, err
		}
		err = sect.MapTo(&c)
	}
	return c, err
}

func writeFile(a *sts.Credentials, c credentials, p string) error {

	wc := &credentials{
		AwsAccessKeyID:  *a.AccessKeyId,
		SecretAccessKey: *a.SecretAccessKey,
		SessionToken:    *a.SessionToken,
		SecurityToken:   *a.SessionToken,
		Region:          c.Region,
		Output:          "json",
		Expiration:      *a.Expiration,
	}

	wc.NamedProfile = c.NamedProfile

	os.OpenFile(credFilePath, os.O_CREATE, 0666)
	cfg, err := ini.Load(credFilePath)
	if err != nil {
		return err
	}
	credSect, err := cfg.NewSection(p)
	err = credSect.ReflectFrom(wc)
	cfg.SaveTo(credFilePath)

	os.OpenFile(configFilePath, os.O_CREATE, 0666)
	cfg, err = ini.Load(configFilePath)
	if err != nil {
		return err
	}
	credSect, err = cfg.NewSection(p)
	err = credSect.ReflectFrom(&credentials{
		Region: c.Region,
		Output: "json",
	})
	cfg.SaveTo(configFilePath)

	return err
}

func validate(arn string, pieces []string) error {
	if strings.Contains(arn, "${") {
		return errors.New("policy variables are not supported")
	}
	if len(pieces) < 6 {
		return errors.New("malformed ARN")
	}
	return nil
}

// Parse accepts and ARN string and attempts to break it into constituent parts.
func parse(arn string) (*components, error) {
	pieces := strings.SplitN(arn, ":", 6)

	if err := validate(arn, pieces); err != nil {
		return nil, err
	}

	components := &components{
		ARN:       pieces[0],
		Partition: pieces[1],
		Service:   pieces[2],
		Region:    pieces[3],
		AccountID: pieces[4],
	}
	if n := strings.Count(pieces[5], ":"); n > 0 {
		components.ResourceDelimiter = ":"
		resourceParts := strings.SplitN(pieces[5], ":", 2)
		components.ResourceType = resourceParts[0]
		components.Resource = resourceParts[1]
	} else {
		if m := strings.Count(pieces[5], "/"); m == 0 {
			components.Resource = pieces[5]
		} else {
			components.ResourceDelimiter = "/"
			resourceParts := strings.SplitN(pieces[5], "/", 2)
			components.ResourceType = resourceParts[0]
			components.Resource = resourceParts[1]
		}
	}
	return components, nil
}

func generateCredentials(c credentials) (*sts.Credentials, error) {

	var err error
	var output interface{}

	s, err := session(c)
	stsSession := sts.New(s)

	callerIdentity, err := stsSession.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	component, err := parse(*callerIdentity.Arn)

	t, err := mfaToken(s, c.MfaSecret)

	if err != nil {
		return nil, err
	}

	switch {
	case c.SourceProfile != "":
		if c.Duration == 0 {
			c.Duration = 3600
		}

		output, err = stsSession.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(c.RoleArn),
			RoleSessionName: aws.String(component.Resource + "-cli"),
			SerialNumber:    aws.String(t.serialNumber),
			TokenCode:       aws.String(t.tokenCode),
			DurationSeconds: aws.Int64(c.Duration),
		})
	case c.SourceProfile == "":
		if c.Duration == 0 {
			c.Duration = 43200
		}
		output, err = stsSession.GetSessionToken(&sts.GetSessionTokenInput{
			DurationSeconds: aws.Int64(c.Duration),
			SerialNumber:    aws.String(t.serialNumber),
			TokenCode:       aws.String(t.tokenCode),
		})
	}

	if a, ok := output.(*sts.AssumeRoleOutput); ok {
		return a.Credentials, err
	}

	if b, ok := output.(*sts.GetSessionTokenOutput); ok {
		return b.Credentials, err
	}

	return nil, errors.New("no session set")
}

func mfaToken(session client.ConfigProvider, secret string) (token, error) {
	devices, err := iam.New(session).ListMFADevices(&iam.ListMFADevicesInput{})
	awsToken := token{}
	if err != nil {
		return awsToken, err
	}
	if len(devices.MFADevices) != 0 {
		token := ""

		if secret == "" {
			token, err = stscreds.StdinTokenProvider()
		} else {
			token, err = totp.GenerateCode(secret, time.Now())
		}
		if err != nil {
			return awsToken, err
		}
		awsToken.tokenCode = token
		awsToken.serialNumber = *devices.MFADevices[0].SerialNumber

	}

	return awsToken, nil
}

func session(c credentials) (*awsSession.Session, error) {
	s := awsSession.Must(awsSession.NewSessionWithOptions(awsSession.Options{
		Config: aws.Config{
			Region: aws.String(c.Region),
			Credentials: awsCredentials.NewCredentials(&awsCredentials.StaticProvider{Value: awsCredentials.Value{
				AccessKeyID:     c.AwsAccessKeyID,
				SecretAccessKey: c.SecretAccessKey,
			}}),
		},
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       awsSession.SharedConfigEnable,
	}))

	if s == nil {
		return s, errors.New("no session set")
	}

	return s, nil
}
