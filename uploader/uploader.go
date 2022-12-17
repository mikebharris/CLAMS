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
	AuthCode     string
	Name         string
	Email        string
	AmountToPay  int
	AmountPaid   int
	DatePaid     string
	Telephone    string
	ArrivalDay   string
	StayingLate  string
	NumberOfKids int
	Diet         string
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

	fmt.Println("Reading from", *csvFile, "and writing to", *sqsQueue)
	for row, record := range records {
		if row == 0 {
			continue // skip row with column headings
		}

		message := Message{
			AuthCode:     record[0],
			Name:         record[1],
			Email:        record[2],
			AmountToPay:  toInt(record[3]),
			AmountPaid:   toInt(record[4]),
			DatePaid:     record[5],
			Telephone:    record[6],
			ArrivalDay:   record[7],
			StayingLate:  record[8],
			NumberOfKids: toInt(record[9]),
			Diet:         record[10],
		}

		if err := sqs.queueMessage(message); err != nil {
			fmt.Println("Unable to queue message ", message, " : ", err)
		} else {
			fmt.Println("Queued message #", row, " : ", message)
		}
	}
}

func toInt(s string) int {
	i, _ := strconv.ParseInt(s, 0, 0)
	return int(i)
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
