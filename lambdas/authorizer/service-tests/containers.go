package main

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"os"
	"time"
)

type Containers struct {
	network         testcontainers.Network
	ssmContainer    testcontainers.Container
	lambdaContainer testcontainers.Container
}

func (c *Containers) Start() error {
	err := c.createNetwork()
	if err != nil {
		return err
	}

	err = c.startSsmContainer()
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

	err := c.ssmContainer.Terminate(context)
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

func (c *Containers) createNetwork() error {
	context := context.Background()
	var err error
	req := testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Driver: "bridge", Name: "myNetwork", Attachable: true},
	}
	c.network, err = testcontainers.GenericNetwork(context, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Containers) startSsmContainer() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	req := testcontainers.ContainerRequest{
		Image:        "wiremock/wiremock",
		ExposedPorts: []string{"8080/tcp"},
		Name:         "ssmmock",
		Hostname:     "ssmmock",
		Networks:     []string{"myNetwork"},
		NetworkMode:  "myNetwork",
		BindMounts:   map[string]string{wd + "/ssm/mappings/": "/home/wiremock/mappings/"},
	}
	context := context.Background()
	c.ssmContainer, err = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return err
	}
	err = c.ssmContainer.Start(context)
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
			"SSM_ENDPOINT_OVERRIDE":   "http://ssmmock:8080",
			"ENVIRONMENT":             "test",
			"DOCKER_LAMBDA_STAY_OPEN": "1",
			"AWS_ACCESS_KEY_ID":       "x",
			"AWS_SECRET_ACCESS_KEY":   "x",
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

func (c *Containers) GetLocalHostLambdaPort() (int, error) {
	context := context.Background()
	mappedPort, err := c.lambdaContainer.MappedPort(context, nat.Port(fmt.Sprintf("%d/tcp", 9001)))
	if err != nil {
		return 0, err
	}
	return mappedPort.Int(), nil
}
