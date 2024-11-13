package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"

	"api/internal/database/db"

	"github.com/rumblefrog/go-a2s"
)

type CsServerPayload struct {
	MaxPlayers    int
	AdminNickname string
}

type CsServerStatusPayload struct {
	// ip:port
	IpAddress string `json:"ipAddress"`
}

type CsService struct {
	db *db.Queries
}

type ServerInfo struct {
	Details   *a2s.ServerInfo
	MapImage  string
	IpAddress string
}

type CsServerStatusResponse struct {
	ServerInfo ServerInfo
	PlayerInfo *a2s.PlayerInfo
	Ping       int
}

type RegisterServerResponse struct {
	IpAddress     string
	AdminNickname string
	AdminPassword string
}

const (
	defaultStartMap = "de_dust2"
)

func NewCsService(db *db.Queries) *CsService {
	return &CsService{
		db: db,
	}
}

func getContainerName(portNumber int) string {
	return fmt.Sprintf("cs_server-%d", portNumber)
}

func createNewCsServer(maxPlayers int, portNumber int, adminNickname, adminPassword string) {
	rootPath := os.Getenv("ROOT_PATH")
	dockerfilePath := fmt.Sprintf("%s/packages/cstrike", rootPath)

	envVars := map[string]string{
		"PORT":           fmt.Sprint(portNumber),
		"MAX_PLAYERS":    fmt.Sprint(maxPlayers),
		"START_MAP":      defaultStartMap,
		"ADMIN_NICKNAME": adminNickname,
		"ADMIN_PASSWORD": adminPassword,
	}

	containerName := getContainerName(portNumber)

	buildCmd := exec.Command(
		"docker",
		"build",
		"--platform",
		"linux/amd64",
		"--build-arg", fmt.Sprintf("%s=%s", "ADMIN_NICKNAME", envVars["ADMIN_NICKNAME"]),
		"--build-arg", fmt.Sprintf("%s=%s", "ADMIN_PASSWORD", envVars["ADMIN_PASSWORD"]),
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

// TODO: Move as a util fn
func generateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *CsService) RegisterServer(ctx context.Context, csServer CsServerPayload) (RegisterServerResponse, error) {
	ports, err := s.db.GetAvailablePorts(ctx)
	if err != nil || len(ports) == 0 {
		return RegisterServerResponse{}, errors.New("no more available ports, please try again later")
	}

	// TODO: Create an algorithm for choosing from available ports
	availablePort := ports[0]

	adminPassword := generateRandomString(10)
	createNewCsServer(csServer.MaxPlayers, int(availablePort.Port), csServer.AdminNickname, adminPassword)

	server, _ := s.db.InsertServer(ctx, db.InsertServerParams{
		MaxPlayers:    sql.NullInt64{Int64: int64(csServer.MaxPlayers), Valid: true},
		AdminNickname: csServer.AdminNickname,
		AdminPassword: adminPassword,
	})

	s.db.UpdatePort(ctx, db.UpdatePortParams{
		ServerID: sql.NullInt64{Int64: int64(server.ID), Valid: true},
		ID:       int64(availablePort.ID),
	})

	machineIpAddress := os.Getenv("IP_ADDRESS")
	return RegisterServerResponse{
		IpAddress:     fmt.Sprintf("%s:%d", machineIpAddress, availablePort.Port),
		AdminNickname: server.AdminNickname,
		AdminPassword: server.AdminPassword,
	}, nil
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

func (s *CsService) DestroyServer(ctx context.Context, serverId int64, port int64) error {
	stopCmd := exec.Command(
		"docker",
		"stop",
		getContainerName(int(port)),
	)

	rmCmd := exec.Command(
		"docker",
		"rm",
		getContainerName(int(port)),
	)

	_, err := stopCmd.Output()
	if err != nil {
		return errors.New("cannot stop server")
	}

	_, err = rmCmd.Output()

	if err != nil {
		return errors.New("cannot remove server")
	}
	portFromDb, _ := s.db.GetPortByValue(ctx, port)

	err = s.db.DeleteServer(ctx, serverId)

	if err != nil {
		return errors.New("cannot delete server from db" + err.Error())
	}
	err = s.db.ResetPort(ctx, portFromDb.ID)

	if err != nil {
		return errors.New("cannot reset port from db")
	}

	return nil
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
		return CsServerStatusResponse{}, fmt.Errorf("error fetching status for server %s", err.Error())
	}

	playerInfo, err := client.QueryPlayer()
	if err != nil {
		return CsServerStatusResponse{}, fmt.Errorf("error fetching status for server %s", err.Error())
	}

	ping, err := pingServer(payload.IpAddress)
	if err != nil {
		return CsServerStatusResponse{}, errors.New("error pinging server")
	}

	mapImage := fmt.Sprintf("https://image.gametracker.com/images/maps/160x120/cs/%s.jpg", serverInfo.Map)
	return CsServerStatusResponse{
		ServerInfo: ServerInfo{
			Details:   serverInfo,
			MapImage:  mapImage,
			IpAddress: payload.IpAddress,
		},
		PlayerInfo: playerInfo,
		Ping:       ping,
	}, nil
}
