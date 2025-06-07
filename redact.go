package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

// RedactorConfig holds configuration for the redactor
type RedactorConfig struct {
	RedactionChar     string
	RedactionPatterns []string
	CompiledPatterns  []*regexp.Regexp
	RedactAll         bool
	ShowPrefix        bool
}

// ConsoleRedactor handles the redaction logic
type ConsoleRedactor struct {
	config *RedactorConfig
}

// NewConsoleRedactor creates a new ConsoleRedactor with default patterns
func NewConsoleRedactor(redactionChar string, redactAll bool, showPrefix bool) *ConsoleRedactor {
	defaultPatterns := []string{
		`\b\d{3}-\d{2}-\d{4}\b`,                               // SSN pattern
		`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`, // Email
		`\b\d{16}\b`,                     // Credit card numbers
		`(?i)\bpassword\s*[:=]\s*\S+`,    // Password patterns
		`(?i)\bapi[_-]?key\s*[:=]\s*\S+`, // API keys
		`(?i)\btoken\s*[:=]\s*\S+`,       // Tokens
		`(?i)\bsecret\s*[:=]\s*\S+`,      // Secrets
		`\b(?:\d{1,3}\.){3}\d{1,3}\b`,    // IP addresses
		`\b[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\b`, // UUIDs
	}

	config := &RedactorConfig{
		RedactionChar:     redactionChar,
		RedactionPatterns: defaultPatterns,
		CompiledPatterns:  make([]*regexp.Regexp, 0),
		RedactAll:         redactAll,
		ShowPrefix:        showPrefix,
	}

	// Compile all patterns
	for _, pattern := range defaultPatterns {
		if compiled, err := regexp.Compile(pattern); err == nil {
			config.CompiledPatterns = append(config.CompiledPatterns, compiled)
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Failed to compile pattern %s: %v\n", pattern, err)
		}
	}

	return &ConsoleRedactor{config: config}
}

// RedactText redacts sensitive information from a string
func (cr *ConsoleRedactor) RedactText(text string) string {
	if cr.config.RedactAll {
		return "[REDACTED]"
	}

	redacted := text
	for _, pattern := range cr.config.CompiledPatterns {
		redacted = pattern.ReplaceAllStringFunc(redacted, func(match string) string {
			return strings.Repeat(cr.config.RedactionChar, len(match))
		})
	}
	return redacted
}

// ProcessReader reads from an io.Reader and processes each line
func (cr *ConsoleRedactor) ProcessReader(reader io.Reader, prefix string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		redactedLine := cr.RedactText(line)

		if cr.config.ShowPrefix && prefix != "" {
			fmt.Printf("%s: %s\n", prefix, redactedLine)
		} else {
			fmt.Println(redactedLine)
		}
	}
}

// RunCommand executes a command and redacts its output
func (cr *ConsoleRedactor) RunCommand(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command specified\n")
		return 1
	}

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
		cr.ProcessReader(stdoutPipe, "STDOUT")
		done <- true
	}()

	// Process stderr
	go func() {
		cr.ProcessReader(stderrPipe, "STDERR")
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

func printUsage() {
	fmt.Fprintf(os.Stderr, `Console Print Redactor

USAGE:
    redactor [OPTIONS] <command> [args...]

OPTIONS:
    -c, --char <CHAR>       Character to use for redaction (default: *)
    -a, --all               Redact all output completely
    -p, --prefix            Show STDOUT/STDERR prefixes
    -h, --help              Show this help message

EXAMPLES:
    redactor python script.py
    redactor -a curl -H "Authorization: Bearer token123" api.example.com
    redactor -c "#" -p node app.js
    redactor go run main.go

DESCRIPTION:
    This tool runs commands and automatically redacts sensitive information
    from their console output, including:

    • Social Security Numbers (XXX-XX-XXXX)
    • Email addresses
    • Credit card numbers
    • API keys and tokens
    • Passwords
    • IP addresses
    • UUIDs

    Use -a/--all flag to completely redact all output for maximum security.
`)
}

func main() {
	var (
		redactionChar = flag.String("c", "*", "Character to use for redaction")
		redactAll     = flag.Bool("a", false, "Redact all output completely")
		showPrefix    = flag.Bool("p", false, "Show STDOUT/STDERR prefixes")
		help          = flag.Bool("h", false, "Show help message")
	)

	// Also support long flags
	flag.StringVar(redactionChar, "char", "*", "Character to use for redaction")
	flag.BoolVar(redactAll, "all", false, "Redact all output completely")
	flag.BoolVar(showPrefix, "prefix", false, "Show STDOUT/STDERR prefixes")
	flag.BoolVar(help, "help", false, "Show help message")

	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command specified\n\n")
		printUsage()
		os.Exit(1)
	}

	// Create redactor
	redactor := NewConsoleRedactor(*redactionChar, *redactAll, *showPrefix)

	// Run command and exit with its exit code
	exitCode := redactor.RunCommand(args)
	os.Exit(exitCode)
}
