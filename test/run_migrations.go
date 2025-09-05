package test

import (
	"os/exec"
	"testing"
)

func RunMigrations(t *testing.T, dsn string) {
	t.Helper()

	cmd := exec.Command("migrate", "-path", "../../migrations", "-database", dsn, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v, output: %s", err, string(output))
	}
	t.Logf("Successfully ran migrations: %s", string(output))
}
