package test

import (
	"context"
{{- if index .Modules "postgres"}}
	"fmt"
{{- end}}
	"net/http/httptest"
	"os"
	"testing"
{{- if index .Modules "postgres"}}
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
{{- end}}

	"github.com/Kroning/example_service/internal/app/config"
	"github.com/Kroning/example_service/internal/app/container"
{{- if index .Modules "postgres"}}
	"github.com/Kroning/example_service/pkg/storage"
	"github.com/Kroning/example_service/pkg/storage/postgresql"
{{- end}}
)

var (
	ts     *httptest.Server
{{- if index .Modules "postgres"}}
	testDB storage.AbstractDB
{{- end}}
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx := context.Background()
{{- if index .Modules "postgres"}}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "test_db",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres container: %v\n", err)
		return 1
	}
	defer pgContainer.Terminate(ctx)

	host, err := pgContainer.Host(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get container host: %v\n", err)
		return 1
	}

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get mapped port: %v\n", err)
		return 1
	}

	testDB = postgresql.New(ctx, postgresql.Config{
		Master: postgresql.ReplicaConfig{
			Host:        host,
			Port:        mappedPort.Port(),
			User:        "test",
			Password:    "test",
			Database:    "test_db",
			MaxOpen:     10,
			MaxIdle:     5,
			MaxLifetime: time.Hour,
			MaxIdleTime: time.Minute,
		},
		Migrations: false,
	})
	defer testDB.Close()

	migrationSQL, err := os.ReadFile("../migrations/db/files/0001_create_table_dummy.up.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read migration file: %v\n", err)
		return 1
	}
	if _, err := testDB.ExecContext(ctx, string(migrationSQL)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run migration: %v\n", err)
		return 1
	}
{{- end}}

	app := container.BuildApp(
		ctx,
		&config.Config{},
{{- if index .Modules "postgres"}}
	testDB,
{{- end}}
	)

	ts = httptest.NewServer(app.Transport.ChiRouter)
	defer ts.Close()

	return m.Run()
}
