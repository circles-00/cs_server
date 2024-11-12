package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"api/internal/database/db"
)

type CsServerPayload struct {
	StartMap   string
	MaxPlayers int
}

type CsService struct {
	db *db.Queries
}

func NewCsService(db *db.Queries) *CsService {
	return &CsService{
		db: db,
	}
}

func createNewCsServer(startMap string, maxPlayers int, portNumber int) {
	dockerfilePath := os.Getenv("DOCKERFILE_PATH")

	envVars := map[string]string{
		"PORT":        fmt.Sprint(portNumber),
		"MAX_PLAYERS": fmt.Sprint(maxPlayers),
		"START_MAP":   startMap,
	}

	containerName := fmt.Sprintf("cs_server-%d", portNumber)

	buildCmd := exec.Command(
		"docker",
		"build",
		"-t",
		containerName,
		dockerfilePath,
	)

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"-p", fmt.Sprintf("%d:%d", portNumber, portNumber),
		"-p", fmt.Sprintf("%d:%d/udp", portNumber, portNumber),
		"--name", containerName,
		"-e", fmt.Sprintf("%s=%s", "PORT", envVars["PORT"]),
		"-e", fmt.Sprintf("%s=%s", "MAX_PLAYERS", envVars["MAX_PLAYERS"]),
		"-e", fmt.Sprintf("%s=%s", "START_MAP", envVars["START_MAP"]),
		containerName,
	)

	out, err := buildCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error building docker image %s %s", string(out), err.Error())
	}

	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error running docker image %s %s", string(out), err.Error())
	}
}

func (s *CsService) RegisterServer(ctx context.Context, csServer CsServerPayload) (int, error) {
	ports, err := s.db.GetAvailablePorts(ctx)
	if err != nil || len(ports) == 0 {
		return -1, errors.New("no more available ports, please try again later")
	}

	// TODO: Create an algorithm for choosing from available ports
	availablePort := ports[0]

	createNewCsServer(csServer.StartMap, csServer.MaxPlayers, int(availablePort.Port))

	server, _ := s.db.InsertServer(ctx, db.InsertServerParams{
		MaxPlayers: sql.NullInt64{Int64: int64(csServer.MaxPlayers), Valid: true},
		StartMap:   sql.NullString{String: csServer.StartMap, Valid: true},
	})

	s.db.UpdatePort(ctx, db.UpdatePortParams{
		ServerID: sql.NullInt64{Int64: int64(server.ID), Valid: true},
		ID:       int64(availablePort.ID),
	})

	return int(availablePort.Port), nil
}
