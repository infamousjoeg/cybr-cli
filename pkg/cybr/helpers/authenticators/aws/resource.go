package aws

// Credential is a struct that holds the AWS credentials
type Credential struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Token           string `json:"Token"`
}

// Resource is used to retrieve the IAM credential. This is required since how an AWS Resource gets an IAM credential depending on its type.AwsResource
// For example lambda gets these credentials via environment variables, while EC2 gets it from a static metadata URL and ECS gets it from a dynamic URL.
type Resource interface {
	Name() string
	GetCredential() (Credential, error)
}
