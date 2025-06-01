package inventory

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// confirm prompts the user for yes/no confirmation and returns true for 'yes' or 'y'.
func confirm(prompt string) bool {
    fmt.Print(prompt)
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(strings.ToLower(input))
    return input == "yes" || input == "y"
}
