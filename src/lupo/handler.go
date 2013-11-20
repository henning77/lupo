package lupo

// Count the connections to make them easily identifiable
var nextConnId = make(chan int)

func genConnectionIds() {
	i := 1
	for {
		nextConnId <- i
		i++
	}
}

func handleConnection(src net.Conn) {
	connId := <-nextConnId
	logPrint(fmt.Sprintf("New connection: %v (from %v)", connId, src.RemoteAddr()))

	var dst net.Conn
	var err error
	if ssl {
		dst, err = tls.Dial("tcp", to, tlsClientConfig())
	} else {
		dst, err = net.Dial("tcp", to)
	}
	
	if err != nil {
		logPrintf("Error connecting to dest: %v", err)
		panic(err)
	}

	// Copy & log in both directions
	go copyWithLog(dst, src, fmt.Sprintf("->%v", connId), "\n")
	go copyWithLog(src, dst, fmt.Sprintf("<-%v", connId), "\n")
}

// Special Writer which writes head + nicely formatted binary / textual chunks + tail
type transferLog struct {
	head string
	tail string
}

// Write nicely formatted binary / textual chunk
func (l *transferLog) Write(p []byte) (n int, err error) {
	logPrint(l.head)

	binChunk := false
	chunkStart := 0
	printableCount := 0

	for i, b := range p {
		if isPrintable(b) {
			printableCount++

			// char chunks have a minimum length
			if binChunk && printableCount >= 5 {
				// Write the previous binary chunk
				chunkEnd := i - printableCount + 1
				writeHexChunk(p[chunkStart:chunkEnd])

				binChunk = false
				chunkStart = chunkEnd
			}
		} else {
			printableCount = 0

			if !binChunk {
				// Write the previous char chunk
				os.Stdout.Write(p[chunkStart:i])

				binChunk = true
				chunkStart = i
			}
		}
	}

	// Final chunk
	if binChunk {
		writeHexChunk(p[chunkStart:])
	} else {
		os.Stdout.Write(p[chunkStart:])
	}

	os.Stdout.WriteString(l.tail)
	return len(p), nil
}

func writeHexChunk(p []byte) {
	dumper := hex.Dumper(os.Stdout)
	dumper.Write(p)
	dumper.Close()
}

func isPrintable(b byte) bool {
	return b == 0x0d || b == 0x0a || (b >= 0x20 && b <= 0x7e)
}

// Copy to dst and to logger
func copyWithLog(dst io.Writer, src io.Reader, head string, tail string) {
	logger := &transferLog{head: head, tail: tail}
	multi := io.MultiWriter(dst, logger)
	if _, err := io.Copy(multi, src); err != nil {
		logPrintf("Closed connection: %v", err)
	}
}
