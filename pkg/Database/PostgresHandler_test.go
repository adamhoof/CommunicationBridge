package database

import (
	env "RPICommandHandler/pkg/Env"
	"fmt"
	"os"
	"testing"
)

func TestPostgresHandler_Connect(t *testing.T) {
	env.SetEnv()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		os.Getenv("dbHost"),
		os.Getenv("dbPort"),
		os.Getenv("dbUser"),
		os.Getenv("dbPassword"),
		os.Getenv("dbName"))

	handler := PostgresHandler{}
	result := handler.Connect(psqlInfo)

	if result != nil {
		t.Errorf("FAILED, expected nil, got %s", result)
		return
	}
	t.Logf("PASSED")
}
