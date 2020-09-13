package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

var awssession = &session.Session{}

func GetSession() *session.Session {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	region := os.Getenv("AWS_REGION")

	session, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secret, sessionToken),
		})
	if err != nil {
		log.Printf("[ERROR] [%s] AWS Session: %s", ServiceName, err)
		return nil
	}
	return session
}

func GetTopic() (string, error) {
	topic, exists := os.LookupEnv("AWS_SNS_TOPIC")
	if !exists {
		log.Printf("[ERROR] [%s] %s", ServiceName, "Topic not set")
		return "", errors.New("topic not set")
	}
	return topic, nil
}
