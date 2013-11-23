package out

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

// Timestamp format for logging
const stamp = "15:04:05.000 "

var Out = os.Stdout

func Stamp(t time.Time) {
	Out.WriteString(t.Format(stamp))
}

func WriteWithoutNewlines(s []byte) {
	// TODO not sure if this is the most efficient way
	scanner := bufio.NewScanner(bytes.NewReader(s))
	for scanner.Scan() {
		Out.Write(scanner.Bytes())
		Out.WriteString(" ")
	}
}

// Log with optimized timestamp
func Printf(s string, a ...interface{}) {
	Out.WriteString(time.Now().Format(stamp))
	Out.WriteString(fmt.Sprintf(s, a...))
	Out.WriteString("\n")
}

// Log with optimized timestamp
func Print(s string) {
	Out.WriteString(time.Now().Format(stamp))
	Out.WriteString(s)
	Out.WriteString("\n")
}
