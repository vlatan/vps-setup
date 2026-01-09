package utils

import (
	"bufio"
	"fmt"
	"strings"
)

// AskQuestion returns user answer
func AskQuestion(prompt string, scanner *bufio.Scanner) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
