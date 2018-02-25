package account

//Account profile
type Account struct {
	ProfileName        string `mapstructure:"profile_name"`
	AwsAccessKeyID     string `mapstructure:"aws_access_key_id"`
	AwsSecretAccessKey string `mapstructure:"aws_secret_access_key"`
	Region             string
	Secret             string
	RoleArn            string `mapstructure:"role_arn"`
	SourceProfile      string `mapstructure:"source_profile"`
}
