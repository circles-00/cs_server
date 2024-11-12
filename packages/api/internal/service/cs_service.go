package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"api/internal/database/db"

	"github.com/rumblefrog/go-a2s"
)

type CsServerPayload struct {
	StartMap   string
	MaxPlayers int
}

type CsServerStatusPayload struct {
	// ip:port
	IpAddress string `json:"ipAddress"`
}

type CsService struct {
	db *db.Queries
}

type ServerInfo struct {
	Details  *a2s.ServerInfo
	MapImage string
}

type CsServerStatusResponse struct {
	ServerInfo ServerInfo
	PlayerInfo *a2s.PlayerInfo
	Ping       int
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

func pingServer(ipAddress string) (int, error) {
	// Note: See https://developer.valvesoftware.com/wiki/Server_queries#A2A_PING
	const a2aPing = "\xFF\xFF\xFF\xFF\x69"

	conn, err := net.Dial("udp", ipAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	totalTime := 0
	tries := 5

	for i := 0; i < tries; i++ {
		start := time.Now()
		_, err = conn.Write([]byte(a2aPing))
		if err != nil {
			return 0, fmt.Errorf("failed to send ping: %w", err)
		}

		buffer := make([]byte, 4096)
		_, err = conn.Read(buffer)
		if err != nil {
			return 0, fmt.Errorf("failed to receive response: %w", err)
		}
		elapsed := time.Since(start)

		totalTime += int(elapsed.Milliseconds())
	}

	averageLatency := totalTime / tries

	return averageLatency, nil
}

func (s *CsService) GetServerList(ctx context.Context) ([]CsServerStatusResponse, error) {
	machineIpAddress := os.Getenv("IP_ADDRESS")
	servers, err := s.db.GetServers(ctx)
	if err != nil {
		return []CsServerStatusResponse{}, errors.New("cannot fetch servers from db")
	}

	serverStatuses := []CsServerStatusResponse{}

	for _, server := range servers {
		status, _ := s.GetServerStatus(ctx, CsServerStatusPayload{IpAddress: fmt.Sprintf("%s:%d", machineIpAddress, server.Port)})
		serverStatuses = append(serverStatuses, status)
	}

	return serverStatuses, nil
}

func (s *CsService) GetServerStatus(ctx context.Context, payload CsServerStatusPayload) (CsServerStatusResponse, error) {
	client, err := a2s.NewClient(payload.IpAddress)
	if err != nil {
		return CsServerStatusResponse{}, errors.New("error fetching status for server")
	}

	serverInfo, err := client.QueryInfo()
	if err != nil {
		return CsServerStatusResponse{}, errors.New("error fetching status for server")
	}

	playerInfo, err := client.QueryPlayer()
	if err != nil {
		return CsServerStatusResponse{}, errors.New("error fetching status for server")
	}

	ping, err := pingServer(payload.IpAddress)
	if err != nil {
		return CsServerStatusResponse{}, errors.New("error pinging server")
	}

	mapImage := fmt.Sprintf("https://image.gametracker.com/images/maps/160x120/cs/%s.jpg", serverInfo.Map)
	return CsServerStatusResponse{
		ServerInfo: ServerInfo{
			Details:  serverInfo,
			MapImage: mapImage,
		},
		PlayerInfo: playerInfo,
		Ping:       ping,
	}, nil
}
