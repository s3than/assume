package main

import (
	"errors"
	"fmt"
	"math/rand"
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

// MfaToken token configuration
type token struct {
	serialNumber string
	tokenCode    string
}

func assumeCommand(args arguments) {

	config, err := getCredentials(args)

	if err != nil {
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
	err = cfg.Section(a).MapTo(&c)

	if !strings.HasPrefix(a, "profile") {
		err = cfg.Section("profile " + a).MapTo(&c)
	}

	if c.SourceProfile != "" {
		err = cfg.Section(c.SourceProfile).MapTo(&c)
		if !strings.HasPrefix(c.SourceProfile, "profile") {
			err = cfg.Section("profile " + c.SourceProfile).MapTo(&c)
		}
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

func generateCredentials(c credentials) (*sts.Credentials, error) {

	var err error
	var output interface{}

	s, err := session(c)
	t, err := mfaToken(s, c.MfaSecret)

	switch {
	case c.SourceProfile != "":
		if c.Duration == 0 {
			c.Duration = 3600
		}
		output, err = sts.New(s).AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(c.RoleArn),
			RoleSessionName: aws.String(randStringBytesMaskSrc(6)),
			SerialNumber:    aws.String(t.serialNumber),
			TokenCode:       aws.String(t.tokenCode),
			DurationSeconds: aws.Int64(c.Duration),
		})
	case c.SourceProfile == "":
		if c.Duration == 0 {
			c.Duration = 43200
		}
		output, err = sts.New(s).GetSessionToken(&sts.GetSessionTokenInput{
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

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// randStringBytesMaskSrc generate random string
func randStringBytesMaskSrc(n int) string {
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
