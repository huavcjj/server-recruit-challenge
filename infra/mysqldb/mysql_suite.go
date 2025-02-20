package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
)

type DBMYSQLSuite struct {
	suite.Suite
	mySQLContainer testcontainers.Container
	ctx            context.Context
	DB             *sql.DB
	user           string
	pass           string
	addr           string
	name           string
}

func (suite *DBMYSQLSuite) SetupTestContainers() (err error) {
	suite.ctx = context.Background()

	req := testcontainers.ContainerRequest{
		Image: "mysql:8.4",
		Env: map[string]string{
			"MYSQL_DATABASE":             "testdb",
			"MYSQL_USER":                 "testuser",
			"MYSQL_PASSWORD":             "testpass",
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server"),
			wait.ForListeningPort("3306/tcp"),
		),
	}

	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	suite.mySQLContainer = container

	host, err := container.Host(suite.ctx)
	if err != nil {
		return err
	}

	port, err := container.MappedPort(suite.ctx, "3306")
	if err != nil {
		return err
	}

	suite.user = "testuser"
	suite.pass = "testpass"
	suite.addr = fmt.Sprintf("%s:%s", host, port.Port()) // 実際のコンテナのポートを使用
	suite.name = "testdb"

	return nil
}

func (suite *DBMYSQLSuite) SetupSuite() {
	err := suite.SetupTestContainers()
	suite.Require().NoError(err)

	suite.DB, err = Initialize(suite.user, suite.pass, suite.addr, suite.name)
	suite.Require().NoError(err)
}

func (suite *DBMYSQLSuite) TearDownSuite() {
	if suite.mySQLContainer != nil {
		err := suite.mySQLContainer.Terminate(suite.ctx)
		suite.Require().NoError(err)
	}
}
