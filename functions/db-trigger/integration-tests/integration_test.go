package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	_ "github.com/lib/pq"
	"github.com/mikebharris/testcontainernetwork-go"
	"github.com/mikebharris/testcontainernetwork-go/clients"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

func TestFeatures(t *testing.T) {
	var steps steps
	steps.t = t
	suite := godog.TestSuite{
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(steps.startContainerNetwork)
			ctx.BeforeSuite(steps.setUpDatabaseClient)
			ctx.AfterSuite(steps.closeDatabaseConnection)
			ctx.AfterSuite(steps.stopContainerNetwork)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^there are database trigger notifications in the database$`, steps.thereAreDatabaseTriggerNotifications)
			ctx.Step(`^the database trigger is invoked$`, steps.theDbTriggerLambdaIsInvoked)
			ctx.Step(`^the messages are placed on the queue$`, steps.theMessagesArePlacedOnTheQueue)
			ctx.Step(`^the notifications are removed from the database$`, steps.theNotificationsAreRemovedFromTheDatabase)
		},
		Options: &godog.Options{
			StopOnFailure: true,
			Strict:        true,
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t, // Testing instance that will run subtests.
		},
	}

	if run := suite.Run(); run != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type steps struct {
	t                         *testing.T
	lambdaContainer           testcontainernetwork.LambdaDockerContainer
	sqsContainer              testcontainernetwork.SqsDockerContainer
	auroraContainer           testcontainernetwork.PostgresDockerContainer
	flywayContainer           testcontainernetwork.FlywayDockerContainer
	networkOfDockerContainers testcontainernetwork.NetworkOfDockerContainers
	dbConx                    *sql.DB
}

func (s *steps) startContainerNetwork() {
	wd, _ := os.Getwd()
	s.sqsContainer = testcontainernetwork.SqsDockerContainer{
		Config: testcontainernetwork.SqsDockerContainerConfig{
			Hostname:       "sqsmock",
			Port:           9324,
			ConfigFilePath: path.Join(wd, "sqs/elasticmq.conf"),
		},
	}

	s.auroraContainer = testcontainernetwork.PostgresDockerContainer{
		Config: testcontainernetwork.PostgresDockerContainerConfig{
			Hostname: "postgres",
			Port:     5432,
			Environment: map[string]string{
				"POSTGRES_PASSWORD": "d0ntHackM3",
				"POSTGRES_USER":     "hacktivista",
				"POSTGRES_DB":       "hacktionlab",
			},
		},
	}

	s.flywayContainer = testcontainernetwork.FlywayDockerContainer{
		Config: testcontainernetwork.FlywayDockerContainerConfig{
			Hostname:       "flyway",
			ConfigFilePath: path.Join(wd, "flyway/conf"),
			SqlFilePath:    path.Join(wd, "../../../flyway/sql"),
		},
	}

	s.lambdaContainer = testcontainernetwork.LambdaDockerContainer{
		Config: testcontainernetwork.LambdaDockerContainerConfig{
			Hostname:   "lambda",
			Executable: "../main",
			Environment: map[string]string{
				"DATABASE_URL": fmt.Sprintf(
					"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
					"postgres", 5432, "hacktivista", "d0ntHackM3", "hacktionlab", "hacktionlab_workshops",
				),
				"DB_TRIGGER_QUEUE":      "db-trigger-queue",
				"SQS_ENDPOINT_OVERRIDE": "http://sqsmock:9324",
			},
		},
	}

	s.networkOfDockerContainers =
		testcontainernetwork.NetworkOfDockerContainers{}.
			WithDockerContainer(&s.lambdaContainer).
			WithDockerContainer(&s.sqsContainer).
			WithDockerContainer(&s.auroraContainer).
			WithDockerContainer(&s.flywayContainer)
	_ = s.networkOfDockerContainers.StartWithDelay(5 * time.Second)
}

func (s *steps) stopContainerNetwork() {
	if err := s.networkOfDockerContainers.Stop(); err != nil {
		log.Fatalf("stopping docker containers: %v", err)
	}
}

func (s *steps) setUpDatabaseClient() {
	var err error
	s.dbConx, err = sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		"localhost", s.auroraContainer.MappedPort(), "hacktivista", "d0ntHackM3", "hacktionlab", "hacktionlab_workshops",
	))
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
}

type TriggerNotification struct {
	WorkshopSignupId int
	WorkshopId       int
	PeopleId         int
	RoleId           int
}

func (a *steps) insertTriggerNotification(n TriggerNotification) {
	msg, _ := json.Marshal(n)
	statement := `insert into trigger_notifications(message) values($1)`
	_, err := a.dbConx.Exec(statement, msg)
	if err != nil {
		panic(err)
	}
}

func (s *steps) thereAreDatabaseTriggerNotifications() {
	s.insertTriggerNotification(TriggerNotification{
		WorkshopSignupId: 1,
		WorkshopId:       2,
		PeopleId:         3,
		RoleId:           4,
	})
	s.insertTriggerNotification(TriggerNotification{
		WorkshopSignupId: 9,
		WorkshopId:       8,
		PeopleId:         7,
		RoleId:           6,
	})
}

func (s *steps) theDbTriggerLambdaIsInvoked() error {
	response, err := http.Post(s.lambdaContainer.InvocationUrl(), "application/json", nil)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.String()
		return fmt.Errorf("invoking Lambda: %d %s", response.StatusCode, body)
	}

	return nil
}

func (s *steps) theMessagesArePlacedOnTheQueue() {
	client, err := clients.SqsClient{}.New(s.sqsContainer.MappedPort())
	if err != nil {
		return
	}

	messages, err := client.GetMessagesFrom("db-trigger-queue")
	if err != nil {
		return
	}

	if len(messages) == 0 {
		panic(fmt.Errorf("expected at least one message in the SQS queue"))
	}

	assert.Len(s.t, messages, 2)

	var tm TriggerNotification
	_ = json.Unmarshal([]byte(*messages[0].Body), &tm)
	assert.Equal(s.t, 1, tm.WorkshopSignupId)
	assert.Equal(s.t, 2, tm.WorkshopId)
	assert.Equal(s.t, 3, tm.PeopleId)
	assert.Equal(s.t, 4, tm.RoleId)

	_ = json.Unmarshal([]byte(*messages[1].Body), &tm)
	assert.Equal(s.t, 9, tm.WorkshopSignupId)
	assert.Equal(s.t, 8, tm.WorkshopId)
	assert.Equal(s.t, 7, tm.PeopleId)
	assert.Equal(s.t, 6, tm.RoleId)
}

func (s *steps) theNotificationsAreRemovedFromTheDatabase() {
	statement := "select count(message) from trigger_notifications"
	row := s.dbConx.QueryRow(statement)
	var count int
	if err := row.Scan(&count); err != nil {
		panic(err)
	}
	assert.Equal(s.t, 0, count)
}

func (s *steps) closeDatabaseConnection() {
	err := s.dbConx.Close()
	if err != nil {
		log.Fatalf("closing db conx: %v", err)
	}
}
