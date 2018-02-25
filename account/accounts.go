package account

//Accounts grouping from config
type Accounts struct {
	Accounts []Account `mapstructure:"accounts"`
}
