package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os/user"
	"time"

	"github.com/fatih/color"
	"github.com/genuinetools/pkg/cli"
	"github.com/go-ini/ini"
	"github.com/s3than/assume/version"
)

var (
	configFile           string
	credFile             string
	profileName          string
	returnProfile           bool
	expiration           bool
	returnNameExpiration bool
	usr, _               = user.Current()
	defaultConfig        = usr.HomeDir + "/.config/assume/config.ini"
	defaultCreds         = usr.HomeDir + "/.config/assume/config.creds"
	configFilePath       = usr.HomeDir + "/.aws/config"
	credFilePath         = usr.HomeDir + "/.aws/credentials"
)

type arguments struct {
	account     string
	saveProfile string
}

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "assume"
	p.Description = "Command line tool to set AWS assume role credentials within the aws credentials files"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("assume", flag.ExitOnError)

	p.FlagSet.StringVar(&configFile, "config", defaultConfig, "config file (default is $HOME/.config/assume/config.ini)")
	p.FlagSet.StringVar(&credFile, "cred", defaultCreds, "credentials file (default is $HOME/.config/assume/config.creds)")
	p.FlagSet.StringVar(&profileName, "p", "default", "set as named profile")
	p.FlagSet.BoolVar(&returnProfile, "d", false, "return name of profile")
	p.FlagSet.BoolVar(&expiration, "t", false, "return expiration time of profile")
	p.FlagSet.BoolVar(&returnNameExpiration, "dt", false, "return expiration time and name of profile")

	p.Action = func(i context.Context, strings []string) error {
		cfg, _ := ini.Load(credFilePath)
		account := "default"
		if len(strings) > 0 {
			account = strings[0]
		}

		switch {
		case returnProfile == false &&
			expiration == false &&
			len(strings) > 0:
			assumeCommand(
				arguments{
					account,
					profileName,
				})
		case expiration != false:
			fmt.Println(remainingTime(cfg, profileName))
		case returnProfile != false:
			fmt.Println(returnProfileName(cfg, profileName))
		case returnNameExpiration != false:
			fmt.Print(returnProfileName(cfg, profileName) + " " + remainingTime(cfg, profileName))
		default:
			return flag.ErrHelp
		}

		return nil
	}

	// Run our program.
	p.Run()
}

func returnProfileName(cfg *ini.File, account string) string {
	sect := cfg.Section(account)
	namedProfile := sect.Key("named_profile").String()

	if !expired(cfg, account) {
		return color.RedString(namedProfile)
	}
	return color.GreenString(namedProfile)
}

func remainingTime(cfg *ini.File, account string) string {
	sect := cfg.Section(account)
	expiration, _ := sect.Key("expiration").Time()
	h, m := fmtDuration(time.Until(expiration))

	if !expired(cfg, account) {
		return color.RedString("%02dh:%02dm", h, m)
	}
	return color.GreenString("%02dh:%02dm", h, m)
}

func expired(cfg *ini.File, account string) bool {
	sect := cfg.Section(account)
	expiration, _ := sect.Key("expiration").Time()
	h, m := fmtDuration(time.Until(expiration))

	if math.Signbit(float64(h)) && math.Signbit(float64(m)) {
		return false
	}
	return true
}

func fmtDuration(d time.Duration) (time.Duration, time.Duration) {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	return h, m
}
