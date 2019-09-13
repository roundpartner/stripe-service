package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"net/http"
)

func (rs *RestServer) WebHook(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s from %s", ServiceName, req.URL.Path, req.RemoteAddr)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(req.Body)
	rs.SNSService.Queue <- buffer

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

type SNSService struct {
	Service *sns.SNS
	Topic   string
	Queue   chan *bytes.Buffer
}

func NewSNSService() *SNSService {
	awssession := GetSession()
	service := sns.New(awssession)
	queue := make(chan *bytes.Buffer, 100)
	topic, _ := GetTopic()
	return &SNSService{
		Service: service,
		Topic:   topic,
		Queue:   queue,
	}
}

func (snsService *SNSService) Run() {
	go func() {
		for {
			buffer := <-snsService.Queue
			err := snsService.Push(buffer)
			if err != nil {
				log.Printf("%s", err.Error())
			}
		}
	}()
}

func (snsService *SNSService) Push(buffer *bytes.Buffer) error {
	params := &sns.PublishInput{
		Subject:  aws.String("Stripe Service Hook"),
		Message:  aws.String(buffer.String()),
		TopicArn: aws.String(snsService.Topic),
	}

	_, err := snsService.Service.Publish(params)

	return err
}
