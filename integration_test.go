package ok_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/broothie/ok/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func Test_integration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.SkipNow()
	}

	for _, tl := range tools.All() {
		t.Run(tl.Name, func(t *testing.T) {
			t.Parallel()

			taskName := fmt.Sprintf("test-%s", strings.ToLower(tl.Name))
			if tl.Name == "Nx" {
				taskName = "integration-test:test-nx"
			}
			taskOutput := fmt.Sprintf("from %s", strings.ToLower(tl.Name))

			t.Run("shows up in tool list", func(t *testing.T) {
				t.Parallel()

				container := startContainer(t)
				defer container.Terminate(context.Background())

				assertCommandOutputContains(t, container, []string{"ok", "--list-tools"}, tl.Name)
			})

			t.Run("shows up in task list", func(t *testing.T) {
				t.Parallel()

				container := startContainer(t)
				defer container.Terminate(context.Background())

				assertCommandOutputContains(t, container, []string{"ok"}, taskName)
			})

			t.Run("run task", func(t *testing.T) {
				t.Parallel()

				container := startContainer(t)
				defer container.Terminate(context.Background())

				assertCommandOutputContains(t, container, []string{"ok", taskName}, taskOutput)
			})
		})
	}
}

func startContainer(t *testing.T) testcontainers.Container {
	t.Helper()

	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:       ".",
				Dockerfile:    "testdata/Dockerfile",
				PrintBuildLog: true,
				KeepImage:     true,
			},
			Cmd: []string{"sleep", "infinity"},
		},
		Started: true,
	})

	require.NoError(t, err)
	return container
}

func assertCommandOutputContains(t *testing.T, container testcontainers.Container, cmd []string, contains string) {
	exitCode, reader, err := container.Exec(context.Background(), cmd)
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)

	outputBytes, err := io.ReadAll(reader)
	require.NoError(t, err)

	assert.Contains(t, string(outputBytes), contains)
}
