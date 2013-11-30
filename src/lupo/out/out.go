package out

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
	"lupo/event"
)

// Timestamp format for logging
const stamp = "2006-01-02 15:04:05.000 "

var Out = os.Stdout

// Generic format:
// <timestamp> <kind><cid> <len> <desc>
func EntryBegin(t time.Time, k event.EventKind, cid event.ConnId, l int) {
	Out.WriteString(t.Format(stamp))
	Out.WriteString(k.String())
	fmt.Fprintf(Out, "%-4d %5d ", cid, l)
}

func EntryEnd() {
	Out.WriteString("\n")
}

func WriteWithoutNewlines(s []byte) {
	// TODO not sure if this is the most efficient way
	scanner := bufio.NewScanner(bytes.NewReader(s))
	for scanner.Scan() {
		Out.Write(scanner.Bytes())
		Out.WriteString(" ")
	}
}

// Directly write an error message and exit
func Fatalf(format string, v ...interface{}) {
	EntryBegin(time.Now(), event.Global, 0, 0)
	fmt.Fprintf(Out, format, v...)
	EntryEnd()
	os.Exit(1)
}