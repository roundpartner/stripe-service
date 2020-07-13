package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"net/http"
	"sync"
)

func (rs *RestServer) WebHook(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Request received: %s %s from %s", ServiceName, req.Method, req.URL.Path, req.RemoteAddr)

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(req.Body)
	if err != nil {
		log.Printf("Error reading buffer: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rs.SNSService.WaitGroup.Add(1)
	rs.SNSService.Queue <- buffer

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

type SNSService struct {
	Service   *sns.SNS
	Topic     string
	Queue     chan *bytes.Buffer
	WaitGroup sync.WaitGroup
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
				log.Printf("Pushing webhook error: %s", err.Error())
			}
			snsService.WaitGroup.Done()
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
