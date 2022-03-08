package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CreateSession creates an AWS session.
func CreateSession(maxRetries int, region string) *session.Session {
	configs := aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(maxRetries),
	}
	return session.Must(session.NewSession(&configs))
}
