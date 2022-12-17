package service_tests

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

type Containers struct {
	network         testcontainers.Network
	dynamoContainer testcontainers.Container
	lambdaContainer testcontainers.Container
}

func (c *Containers) Start() error {
	err := c.createNetwork()
	if err != nil {
		return err
	}

	err = c.startDynamoDbContainer()
	if err != nil {
		return err
	}

	err = c.startLambdaContainer()
	if err != nil {
		return err
	}

	fmt.Println("Sleeping for 5 seconds while containers start")
	time.Sleep(5 * time.Second)
	return nil
}

func (c *Containers) Stop() error {
	context := context.Background()

	err := c.dynamoContainer.Terminate(context)
	if err != nil {
		return err
	}

	err = c.lambdaContainer.Terminate(context)
	if err != nil {
		return err
	}

	err = c.network.Remove(context)
	if err != nil {
		return err
	}
	return nil
}

func (c *Containers) GetLambdaLog() io.ReadCloser {
	logs, err := c.lambdaContainer.Logs(context.Background())
	if err != nil {
		panic(err)
	}
	return logs
}

func (c *Containers) GetLocalhostPort(container testcontainers.Container, port int) int {
	mappedPort, err := container.MappedPort(context.Background(), nat.Port(fmt.Sprintf("%d/tcp", port)))
	if err != nil {
		panic(err)
	}
	return mappedPort.Int()
}

func (c *Containers) GetLocalHostDynamoPort() int {
	return c.GetLocalhostPort(c.dynamoContainer, 8000)
}

func (c *Containers) GetLocalHostLambdaPort() int {
	return c.GetLocalhostPort(c.lambdaContainer, 9001)
}

func (c *Containers) createNetwork() error {
	var err error
	req := testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Driver: "bridge", Name: "myNetwork", Attachable: true},
	}
	c.network, err = testcontainers.GenericNetwork(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Containers) startDynamoDbContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
		Name:         "dynamo",
		Hostname:     "dynamo",
		Networks:     []string{"myNetwork"},
		NetworkMode:  "myNetwork",
		Entrypoint:   []string{"java", "-jar", "DynamoDBLocal.jar", "-inMemory", "-sharedDb"},
	}
	var err error
	context := context.Background()
	c.dynamoContainer, err = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	err = c.dynamoContainer.Start(context)
	if err != nil {
		return err
	}

	return nil
}

func (c *Containers) startLambdaContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "lambci/lambda:go1.x",
		ExposedPorts: []string{"9001/tcp"},
		Name:         "lambda",
		Hostname:     "lambda",
		Env: map[string]string{
			"ATTENDEES_TABLE_NAME":     "attendees",
			"DYNAMO_ENDPOINT_OVERRIDE": "http://dynamo:8000",
			"DOCKER_LAMBDA_STAY_OPEN":  "1",
			"AWS_ACCESS_KEY_ID":        "x",
			"AWS_SECRET_ACCESS_KEY":    "x",
		},
		Networks:    []string{"myNetwork"},
		NetworkMode: "myNetwork",
	}
	context := context.Background()
	var err error
	c.lambdaContainer, err = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return err
	}
	err = c.lambdaContainer.CopyFileToContainer(context, "main", "/var/task/handler", 365)
	if err != nil {
		return err
	}
	c.lambdaContainer.Start(context)

	return nil
}
