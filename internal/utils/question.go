package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// AskQuestion returns user answer from terminal
func AskQuestion(prompt string, scanner *bufio.Scanner) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// AskSensitiveQuestion returns user sensitive answer from terminal
func AskSensitiveQuestion(prompt string) (string, error) {
	fmt.Print(prompt)
	value, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(value), nil
}
