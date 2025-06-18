package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
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
		// Forward signal to the child process
		if cmd.Process != nil {
			cmd.Process.Signal(os.Interrupt)
			// Give the process a chance to handle the signal gracefully
			time.Sleep(100 * time.Millisecond)
			// If still running, force kill
			if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
				cmd.Process.Kill()
			}
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
		// Wait for output processing to finish with timeout
		select {
		case <-outputDone:
		case <-time.After(1 * time.Second):
			// Timeout waiting for output processing
		}
		cancel()
		return exitCode
	case <-ctx.Done():
		// Wait for command to finish with timeout
		select {
		case exitCode := <-cmdDone:
			return exitCode
		case <-time.After(2 * time.Second):
			// Force kill if still running
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			return 1
		}
	}
}

func (r *redactor) handleStdin(ctx context.Context, ptyMaster pty.Pty) {
	// Use a goroutine to handle stdin reading without blocking
	stdinChan := make(chan []byte, 1)
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				}
				return
			}
			if n > 0 {
				// Make a copy of the buffer to send through channel
				data := make([]byte, n)
				copy(data, buffer[:n])
				select {
				case stdinChan <- data:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-stdinChan:
			_, writeErr := ptyMaster.Write(data)
			if writeErr != nil {
				fmt.Fprintf(os.Stderr, "Error writing to PTY: %v\n", writeErr)
				return
			}
		}
	}
}

func (r *redactor) handleOutput(ctx context.Context, ptyMaster pty.Pty, done chan<- int) {
	buffer := make([]byte, 4096)
	defer func() {
		done <- 0
	}()

	// Use a goroutine to handle PTY reading without blocking
	outputChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)

	go func() {
		for {
			n, err := ptyMaster.Read(buffer)
			if n > 0 {
				// Make a copy of the buffer to send through channel
				data := make([]byte, n)
				copy(data, buffer[:n])
				select {
				case outputChan <- data:
				case <-ctx.Done():
					return
				}
			}
			if err != nil {
				select {
				case errorChan <- err:
				case <-ctx.Done():
				}
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-outputChan:
			// Process and redact the output
			text := string(data)
			redactedText := r.processAndRedactText(text)

			// Write the redacted content to stdout
			_, writeErr := os.Stdout.Write([]byte(redactedText))
			if writeErr != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", writeErr)
				return
			}
		case err := <-errorChan:
			if err == io.EOF {
				// Command finished
				return
			}
			fmt.Fprintf(os.Stderr, "Error reading PTY: %v\n", err)
			return
		}
	}
}

func (r *redactor) processAndRedactText(text string) string {
	if len(text) == 0 {
		return text
	}

	// Use comprehensive token-based redaction
	return r.redactTokens(text)
}

// redactTokens implements comprehensive token-based redaction
func (r *redactor) redactTokens(text string) string {
	if len(r.redactText) == 0 {
		return text
	}

	redactedText := text

	// Split text into tokens using multiple delimiters
	tokens := r.tokenizeText(text)

	// Process each token
	for _, token := range tokens {
		if token == "" {
			continue
		}

		// Check if this token contains any sensitive value
		shouldRedact := r.tokenContainsSensitiveData(token)

		if shouldRedact {
			// Replace the entire token in the text
			replacement := r.generateRedaction(len(token))
			redactedText = strings.ReplaceAll(redactedText, token, replacement)
		}
	}

	return redactedText
}

// tokenizeText splits text into meaningful tokens using various delimiters
func (r *redactor) tokenizeText(text string) []string {
	var tokens []string

	// First split by whitespace to get main chunks
	fields := strings.Fields(text)

	for _, field := range fields {
		// Further split each field by common delimiters while preserving the original tokens
		subTokens := r.splitByDelimiters(field)
		tokens = append(tokens, subTokens...)
	}

	return tokens
}

// splitByDelimiters splits a string by various delimiters but keeps meaningful tokens together
func (r *redactor) splitByDelimiters(s string) []string {
	if s == "" {
		return nil
	}

	var tokens []string
	currentToken := ""

	// Common delimiters that might separate tokens
	delimiters := " \t\n\r,;|&"

	for i, char := range s {
		if strings.ContainsRune(delimiters, char) {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
		} else {
			currentToken += string(char)
		}

		// If we're at the end, add the current token
		if i == len(s)-1 && currentToken != "" {
			tokens = append(tokens, currentToken)
		}
	}

	// Also add the original string as a token to catch cases where
	// sensitive data spans across what we consider delimiters
	tokens = append(tokens, s)

	return tokens
}

// tokenContainsSensitiveData checks if a token contains any sensitive data
func (r *redactor) tokenContainsSensitiveData(token string) bool {
	for _, valueToRedact := range r.redactText {
		if valueToRedact == "" {
			continue
		}

		// Check direct containment
		if strings.Contains(token, valueToRedact) {
			return true
		}

		// Check after cleaning the token
		cleanToken := r.cleanWordForMatching(token)
		if strings.Contains(cleanToken, valueToRedact) {
			return true
		}

		// Check if the sensitive value is contained in the cleaned token
		cleanValue := r.cleanWordForMatching(valueToRedact)
		if cleanValue != "" && strings.Contains(cleanToken, cleanValue) {
			return true
		}
	}

	return false
}

func (r *redactor) cleanWordForMatching(word string) string {
	// Remove common punctuation and special characters that might be attached to sensitive values
	cleaned := strings.TrimFunc(word, func(r rune) bool {
		return r == ',' || r == '.' || r == ';' || r == ':' || r == '"' || r == '\'' ||
			r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' ||
			r == '\n' || r == '\r' || r == '\t' || r == '@' || r == '/' || r == '?' ||
			r == '&' || r == '=' || r == '#' || r == '%' || r == '+' || r == '~'
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
