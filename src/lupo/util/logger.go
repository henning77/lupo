package util

import (
	"time"
	"os"
	"fmt"
)

// Timestamp format for logging
const stamp = "15:04:05.000"

// Log with optimized timestamp
func Printf(s string, a ...interface{}) {
	os.Stdout.WriteString(time.Now().Format(stamp))
	os.Stdout.WriteString(" ")
	os.Stdout.WriteString(fmt.Sprintf(s, a...))
	os.Stdout.WriteString("\n")
}

// Log with optimized timestamp
func Print(s string) {
	os.Stdout.WriteString(time.Now().Format(stamp))
	os.Stdout.WriteString(" ")
	os.Stdout.WriteString(s)
	os.Stdout.WriteString("\n")
}
