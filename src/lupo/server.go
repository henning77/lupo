package lupo

// Host:port to listen from
var From string

// Host:port to forward to
var To string

// If true, use SSL/TLS connections
var Ssl bool


func tlsServerConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		logPrintf("Could not read server certificate (cert.pem, key.pem): %v", err)
		os.Exit(1)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func tlsClientConfig() *tls.Config {
	// Simply accept everything
	return &tls.Config{InsecureSkipVerify: true}
}


func Listen() {
	Printf("Listening to [%v], forwarding to [%v]", From, To)

	var ln net.Listener
	var err error
	if Ssl {
		ln, err = tls.Listen("tcp", From, tlsServerConfig())
	} else {
		ln, err = net.Listen("tcp", From)
	}
	
	if err != nil {
		logPrintf("Could not open port: %v", err)
		os.Exit(1)
	}

	go genConnectionIds()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logPrintf("Error accepting: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}