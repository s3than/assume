package account

//Account profile
type Account struct {
	ProfileName        string `mapstructure:"profile_name" yaml:"profile_name,omitempty"`
	AwsAccessKeyID     string `mapstructure:"aws_access_key_id" yaml:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey string `mapstructure:"aws_secret_access_key" yaml:"aws_secret_access_key,omitempty"`
	Region             string `yaml:"region,omitempty"`
	Secret             string `yaml:"secret,omitempty"`
	RoleArn            string `mapstructure:"role_arn" yaml:"role_arn,omitempty"`
	SourceProfile      string `mapstructure:"source_profile" yaml:"source_profile,omitempty"`
}

func (a Account) IsBase() bool {
	if a.SourceProfile == "" && a.AwsAccessKeyID != "" {
		return true
	}
	return false
}
