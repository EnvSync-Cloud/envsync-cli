package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/aymanbagabas/go-pty"
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
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command provided\n")
		return 1
	}

	// Create a new PTY
	ptyMaster, err := pty.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PTY: %v\n", err)
		return 1
	}
	defer ptyMaster.Close()

	// Create the command using PTY
	cmd := ptyMaster.Command(args[0], args[1:]...)

	// Set environment variables to force color output
	cmd.Env = append(os.Environ(),
		"FORCE_COLOR=1",
		"CLICOLOR_FORCE=1",
		"TERM=xterm-256color",
	)

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting command: %v\n", err)
		return 1
	}

	// Handle interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		<-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived interrupt signal, terminating command...\n")
		cancel()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Channels to signal completion
	outputDone := make(chan int, 1)
	cmdDone := make(chan int, 1)

	// Handle stdin in a separate goroutine
	go r.handleStdin(ctx, ptyMaster)

	// Handle stdout/stderr processing
	go r.handleOutput(ctx, ptyMaster, outputDone)

	// Wait for command completion or context cancellation
	go func() {
		exitCode := 0
		if err := cmd.Wait(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = 1
			}
		}
		cmdDone <- exitCode
	}()

	// Wait for completion or cancellation
	select {
	case exitCode := <-cmdDone:
		// Wait for output processing to finish
		<-outputDone
		cancel()
		return exitCode
	case <-ctx.Done():
		return 1
	}
}

func (r *redactor) handleStdin(ctx context.Context, ptyMaster pty.Pty) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Set a short read timeout to allow checking context
			os.Stdin.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, err := os.Stdin.Read(buffer)
			if err != nil {
				if !isTimeoutError(err) && err != io.EOF {
					fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				}
				continue
			}

			if n > 0 {
				_, writeErr := ptyMaster.Write(buffer[:n])
				if writeErr != nil {
					fmt.Fprintf(os.Stderr, "Error writing to PTY: %v\n", writeErr)
					return
				}
			}
		}
	}
}

func (r *redactor) handleOutput(ctx context.Context, ptyMaster pty.Pty, done chan<- int) {
	buffer := make([]byte, 4096)
	defer func() {
		done <- 0
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := ptyMaster.Read(buffer)
			if n > 0 {
				// Process and redact the output
				text := string(buffer[:n])
				redactedText := r.processAndRedactText(text)

				// Write the redacted content to stdout
				_, writeErr := os.Stdout.Write([]byte(redactedText))
				if writeErr != nil {
					fmt.Fprintf(os.Stderr, "Error writing output: %v\n", writeErr)
					return
				}
			}

			if err != nil {
				if err == io.EOF {
					// Command finished
					return
				}
				fmt.Fprintf(os.Stderr, "Error reading PTY: %v\n", err)
				return
			}
		}
	}
}

func (r *redactor) processAndRedactText(text string) string {
	// Split text preserving whitespace and special characters
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	// Create a map to track which words to redact
	redactedText := text

	// Process each word for redaction
	for _, word := range words {
		if word == "" {
			continue
		}

		// Clean the word of special characters for matching
		cleanWord := r.cleanWordForMatching(word)

		// Check if this word should be redacted
		if slices.Contains(r.redactText, cleanWord) {
			// Replace the word in the original text
			redactedText = strings.ReplaceAll(redactedText, word, r.generateRedaction(len(cleanWord)))
		}
	}

	return redactedText
}

func (r *redactor) cleanWordForMatching(word string) string {
	// Remove common punctuation and special characters that might be attached to sensitive values
	cleaned := strings.TrimFunc(word, func(r rune) bool {
		return r == ',' || r == '.' || r == ';' || r == ':' || r == '"' || r == '\'' ||
			r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' ||
			r == '\n' || r == '\r' || r == '\t'
	})
	return cleaned
}

func (r *redactor) generateRedaction(length int) string {
	if length <= 0 {
		return "********"
	}
	// Generate redaction based on original length, but minimum 8 characters
	redactionLength := length
	if redactionLength < 8 {
		redactionLength = 8
	}
	return strings.Repeat("*", redactionLength)
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	// Check if it's a timeout error
	if netErr, ok := err.(interface{ Timeout() bool }); ok {
		return netErr.Timeout()
	}
	return false
}
