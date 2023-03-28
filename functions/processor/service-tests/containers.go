package service_tests

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"os"
	"path/filepath"
)

type Containers struct {
	network         testcontainers.Network
	auroraContainer testcontainers.Container
	dynamoContainer testcontainers.Container
	flywayContainer testcontainers.Container
	lambdaContainer testcontainers.Container
}

func (c *Containers) start() {
	c.createNetwork()
	c.startDynamoDbContainer()
	c.startAuroraContainer()
	c.startFlywayContainer()
	c.startLambdaContainer()
}

func (c *Containers) stop() {
	ctx := context.Background()

	if err := c.lambdaContainer.Terminate(ctx); err != nil {
		panic(err)
	}

	if err := c.auroraContainer.Terminate(ctx); err != nil {
		panic(err)
	}

	if err := c.dynamoContainer.Terminate(ctx); err != nil {
		panic(err)
	}

	if err := c.network.Remove(ctx); err != nil {
		panic(err)
	}
}

func (c *Containers) getLambdaLog() io.ReadCloser {
	logs, err := c.lambdaContainer.Logs(context.Background())
	if err != nil {
		panic(err)
	}
	return logs
}

func (c *Containers) getLocalhostPort(container testcontainers.Container, port int) int {
	mappedPort, err := container.MappedPort(context.Background(), nat.Port(fmt.Sprintf("%d/tcp", port)))
	if err != nil {
		panic(err)
	}
	return mappedPort.Int()
}

func (c *Containers) getLambdaPort() int {
	return c.getLocalhostPort(c.lambdaContainer, 9001)
}

func (c *Containers) getAuroraPort() int {
	return c.getLocalhostPort(c.auroraContainer, 5432)
}

func (c *Containers) getDynamoPort() int {
	return c.getLocalhostPort(c.dynamoContainer, 8000)
}

func (c *Containers) createNetwork() {
	var err error
	req := testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Driver: "bridge", Name: "myNetwork", Attachable: true},
	}
	c.network, err = testcontainers.GenericNetwork(context.Background(), req)
	if err != nil {
		panic(err)
	}
}

func (c *Containers) startLambdaContainer() {
	req := testcontainers.ContainerRequest{
		Image:        "lambci/lambda:go1.x",
		ExposedPorts: []string{"9001/tcp"},
		Name:         "lambda",
		Hostname:     "lambda",
		Env: map[string]string{
			"DOCKER_LAMBDA_STAY_OPEN":     "1",
			"AWS_ACCESS_KEY_ID":           "x",
			"AWS_SECRET_ACCESS_KEY":       "x",
			"ENVIRONMENT":                 "test",
			"DB_HOST":                     "postgres",
			"DB_NAME":                     "hacktionlab",
			"DB_USER":                     "hacktivista",
			"DB_PASSWORD":                 "d0ntHackM3",
			"WORKSHOP_SIGNUPS_TABLE_NAME": workshopSignupsTableName,
			"DYNAMO_ENDPOINT_OVERRIDE":    "http://dynamo:8000",
		},
		Networks:    []string{"myNetwork"},
		NetworkMode: "myNetwork",
	}
	ctx := context.Background()
	var err error
	c.lambdaContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		panic(err)
	}
	if err := c.lambdaContainer.CopyFileToContainer(ctx, "main", "/var/task/handler", 365); err != nil {
		panic(err)
	}
	c.lambdaContainer.Start(ctx)
}

func (c *Containers) startAuroraContainer() {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Name:         "postgres",
		Hostname:     "postgres",
		Networks:     []string{"myNetwork"},
		NetworkMode:  "myNetwork",
		Env: map[string]string{
			"POSTGRES_PASSWORD": "d0ntHackM3",
			"POSTGRES_USER":     "hacktivista",
			"POSTGRES_DB":       "hacktionlab",
		},
	}

	ctx := context.Background()

	c.auroraContainer, _ = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})

	if err := c.auroraContainer.Start(ctx); err != nil {
		panic(err)
	}
}

func (c *Containers) startFlywayContainer() {
	cwd, _ := os.Getwd()

	req := testcontainers.ContainerRequest{
		Image:       "flyway/flyway",
		Name:        "flyway",
		Hostname:    "flyway",
		Networks:    []string{"myNetwork"},
		NetworkMode: "myNetwork",
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source:   testcontainers.GenericBindMountSource{HostPath: filepath.Join(cwd, "..", "..", "..", "flyway", "sql")},
				Target:   "/flyway/sql",
				ReadOnly: true,
			},
			testcontainers.ContainerMount{
				Source:   testcontainers.GenericBindMountSource{HostPath: filepath.Join(cwd, "flyway", "conf")},
				Target:   "/flyway/conf",
				ReadOnly: true,
			},
		},
		Entrypoint: []string{"flyway", "migrate"},
	}

	ctx := context.Background()

	var err error
	c.flywayContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		panic(err)
	}

	if err := c.flywayContainer.Start(ctx); err != nil {
		panic(err)
	}
}

func (c *Containers) startDynamoDbContainer() {
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
	ctx := context.Background()
	c.dynamoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		panic(err)
	}
	if err := c.dynamoContainer.Start(ctx); err != nil {
		panic(err)
	}
}
