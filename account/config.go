package account

import (
	"errors"
	"reflect"

	"github.com/spf13/viper"
)

func mapAccountsByField(field string) map[string]Account {
	var config Accounts

	viper.Unmarshal(&config)

	confMap := map[string]Account{}
	for _, v := range config.Accounts {
		r := reflect.ValueOf(v)

		if r.FieldByName(field).IsValid() == true {
			f := reflect.Indirect(r).FieldByName(field)
			value := f.Interface().(string)
			confMap[value] = v
		}
	}
	return confMap
}

// WriteAccountToConfig write account to config file
func WriteAccountToConfig(account Account) bool {

	var config Accounts

	viper.Unmarshal(&config)
	newAccounts := append(config.Accounts, account)

	viper.Set("accounts", newAccounts)
	err := viper.WriteConfig()

	if err != nil {
		return false
	}
	return true
}

func ProfileExists(profileName string) bool {

	confMap := mapAccountsByField("ProfileName")

	if _, ok := confMap[profileName]; ok {
		return true
	}

	return false
}

// FindAllbyType return accounts by type
func FindAllbyType(accountType string) ([]Account, error) {

	confMap := mapAccountsByField("ProfileName")

	if allowedTypes(accountType) == false {
		return nil, errors.New("invalid account type")
	}

	var accounts []Account

	for _, a := range confMap {
		if a.IsBase() && accountType == "base" {
			accounts = append(accounts, a)
		} else if !a.IsBase() && accountType == "cross" {
			accounts = append(accounts, a)
		}
	}

	return accounts, nil
}

func allowedTypes(accountType string) bool {
	if accountType == "base" || accountType == "cross" {
		return true
	}

	return false
}
