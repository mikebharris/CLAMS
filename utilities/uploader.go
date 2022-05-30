package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"os"
	"strconv"
)

type Message struct {
	Name        string
	Email       string
	Code        string
	ToPay       uint
	Paid        uint
	PaidDate    string
	Phone       string
	Arrival     string
	Diet        string
	StayingLate string
	Kids        uint
}

var csvFile = flag.String("csv", "", "input csv file name")
var sqsQueue = flag.String("sqs", "", "output sqs queue")

func main() {
	flag.Parse()

	f, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Printf("reading from CSV file: %v", err)
	}

	if len(records) == 0 {
		log.Println("no records to upload")
		return
	}

	sqs, err := newSqsClient()

	for row, record := range records {
		if row == 0 {
			continue
		}

		message := Message{
			Name:        record[0],
			Email:       record[1],
			Code:        record[2],
			ToPay:       toUint(record[3]),
			Paid:        toUint(record[4]),
			PaidDate:    record[5],
			Phone:       record[6],
			Arrival:     record[7],
			Diet:        record[8],
			StayingLate: record[9],
			Kids:        toUint(record[10]),
		}

		if err := sqs.queueMessage(message); err != nil {
			fmt.Println("unable to queue message ", message, " : ", err)
		} else {
			fmt.Println("successfully queued message #", row, " : ", message)
		}
	}
}

func toUint(s string) uint {
	u, _ := strconv.ParseUint(s, 0, 0)
	return uint(u)
}

type SqsClient struct {
	sqsHandle *sqs.Client
}

func newSqsClient() (SqsClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	client := sqs.NewFromConfig(cfg)

	return SqsClient{sqsHandle: client}, nil
}

func (s *SqsClient) getQueueUrl() *string {
	output, err := s.sqsHandle.GetQueueUrl(
		context.Background(),
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(*sqsQueue),
		},
	)
	if err != nil {
		panic(err)
	}

	return output.QueueUrl
}

func (s *SqsClient) queueMessage(message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	s.sqsHandle.SendMessage(
		context.Background(),
		&sqs.SendMessageInput{
			MessageBody: aws.String(string(body)),
			QueueUrl:    s.getQueueUrl(),
		},
	)

	return nil
}
