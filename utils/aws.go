package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

// GetAWSSession creates a session and ensures credentials are present
func GetAWSSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	// Check if we have credentials
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return nil, fmt.Errorf("AWS credentials not found. Please run 'gosak assume' or configure your credentials.")
	}

	return sess, nil
}
