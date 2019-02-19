package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/client"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
)

func Test_assumeCommand(t *testing.T) {
	type args struct {
		args arguments
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assumeCommand(tt.args.args)
		})
	}
}

func Test_getSection(t *testing.T) {
	type args struct {
		cfg *ini.File
		a   string
	}
	tests := []struct {
		name    string
		args    args
		want    *ini.Section
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSection(tt.args.cfg, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCredentials(t *testing.T) {
	type args struct {
		args arguments
	}
	tests := []struct {
		name    string
		args    args
		want    credentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCredentials(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeFile(t *testing.T) {
	type args struct {
		a *sts.Credentials
		c credentials
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeFile(tt.args.a, tt.args.c, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		arn    string
		pieces []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.arn, tt.args.pieces); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		arn string
	}
	tests := []struct {
		name    string
		args    args
		want    *components
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.arn)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCredentials(t *testing.T) {
	type args struct {
		c credentials
	}
	tests := []struct {
		name    string
		args    args
		want    *sts.Credentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateCredentials(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mfaToken(t *testing.T) {
	type args struct {
		session client.ConfigProvider
		secret  string
	}
	tests := []struct {
		name    string
		args    args
		want    token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mfaToken(tt.args.session, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("mfaToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mfaToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_session(t *testing.T) {
	type args struct {
		c credentials
	}
	tests := []struct {
		name    string
		args    args
		want    *awsSession.Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := session(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("session() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("session() = %v, want %v", got, tt.want)
			}
		})
	}
}
