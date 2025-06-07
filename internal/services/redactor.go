package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type RedactorService interface {
	RunRedactor(args []string) int
}

type redactor struct {
	redactText []string
}

func NewRedactorService(redactText []string) RedactorService {
	return &redactor{
		redactText: redactText,
	}
}

func (r *redactor) RunRedactor(args []string) int {
	// Create command
	cmd := exec.Command(args[0], args[1:]...)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		return 1
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stderr pipe: %v\n", err)
		return 1
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting command: %v\n", err)
		return 1
	}

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Process output concurrently
	done := make(chan bool, 2)

	// Process stdout
	go func() {
		r.ProcessReader(stdoutPipe)
		done <- true
	}()

	// Process stderr
	go func() {
		r.ProcessReader(stderrPipe)
		done <- true
	}()

	// Wait for command completion or interruption
	go func() {
		<-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived interrupt signal, terminating command...\n")
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Wait for both readers to finish
	<-done
	<-done

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		return 1
	}

	return 0
}

func (r *redactor) ProcessReader(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		redactedLine := r.RedactText(line)
		fmt.Printf("%s", redactedLine)
	}
}

func (r *redactor) RedactText(text string) string {
	words := strings.Fields(text)

	for i, word := range words {
		for _, redactText := range r.redactText {
			if word == redactText {
				words[i] = strings.Repeat("*", 8)
			}
		}
	}

	return strings.Join(words, " ")
}
