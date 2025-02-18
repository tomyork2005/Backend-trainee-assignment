package storage

import (
	"Backend-trainee-assignment/internal/model/service"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func startTestContainer(t *testing.T) string {
	ctx := context.Background()

	absPath1, err := filepath.Abs("../../migrations/0001_init.up.sql")
	absPath2, err := filepath.Abs("../../migrations/0002_seed_data.sql")
	require.NoError(t, err)

	r1, err := os.Open(absPath1)
	r2, err := os.Open(absPath2)
	require.NoError(t, err)

	req := testcontainers.ContainerRequest{
		Image: "postgres:latest",
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            r1,
				ContainerFilePath: "/docker-entrypoint-initdb.d/0001_init.up.sql",
				FileMode:          0o644,
			},
			{
				Reader:            r2,
				ContainerFilePath: "/docker-entrypoint-initdb.d/0002_seed_data.sql",
				FileMode:          0o644,
			},
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	testcontainers.CleanupContainer(t, postgresC)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port())

	return connStr
}

func Test_GetUserByUsername(t *testing.T) {

	// Arrange

	tests := []struct {
		name     string
		username string

		wantErr      error
		expectedUser *service.User
	}{
		{
			name:     "happy_path",
			username: "testUserExist",
			wantErr:  nil,
			expectedUser: &service.User{
				Username: "testUserExist",
				Password: "testPassExist",
				Balance:  1000,
			},
		},
		{
			name:         "user not found",
			username:     "testUserNotExist",
			wantErr:      nil,
			expectedUser: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			connStr := startTestContainer(t)
			storage, err := NewStorage(connStr)
			require.NoError(t, err)

			// Act

			user, err := storage.GetUserByUsername(context.Background(), tt.username)

			// Assert
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.expectedUser, user)

		})
	}
}
