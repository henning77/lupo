lupo. simple network logging.
=============================

lupo is useful for debugging network communications.

lupo is tiny & easy to use. Use it when you don't need or want a full blown tcpdump/wireshark.

lupo listens to a tcp port, and forwards it to another.
All communication which goes through, both ways, is logged to stdout in a convenient format.

Features:

* No installation needed. lupo runs as a standalone executable.
* Cross-platform. Binary available for all major platforms.
* Works with SSL. Lets you see the encrypted traffic in plain text.
* Choose between full or short output format with the `-style` flag.
* Log statistics to csv (`-style stats`).

Get the latest binary 
---------------------

[Windows](https://github.com/henning77/lupo/releases/download/v0.2.0/lupo_0.2.0_windows_amd64.zip) |
[Mac OS X](https://github.com/henning77/lupo/releases/download/v0.2.0/lupo_0.2.0_darwin_amd64.zip) |
[Linux tarball](https://github.com/henning77/lupo/releases/download/v0.2.0/lupo_0.2.0_linux_amd64.tar.gz) |
[Linux .deb](https://github.com/henning77/lupo/releases/download/v0.2.0/lupo_0.2.0_amd64.deb)

Example 1: Log http traffic
---------------------------

	lupo -from :8080 -to google.com:80

Point your browser to http://localhost:8080. lupo will print something like:

	2013-11-30 14:36:26.836   0   0 Listening to [:8080], forwarding to [google.com:80]
	2013-11-30 14:36:50.068  [1   0 New connection from [::1]:65468
	2013-11-30 14:36:50.069 ->1 360 GET / HTTP/1.1
	2013-11-30 14:36:50.105 <-1 506 HTTP/1.1 302 Found <HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8"> < (...)
	2013-11-30 14:38:18.213  ]1   0 Client closed connection

Example 2: Full logging
-----------------------

	If you prefer to see everything transmitted, use `-style full`.

	lupo -from :8080 -to google.com:80 -style full

	2013-11-30 14:36:26.836   0   0 Listening to [:8080], forwarding to [google.com:80]
	2013-11-30 14:36:50.068  [1   0 New connection from [::1]:65468
	2013-11-30 14:36:50.069 ->1 360 GET / HTTP/1.1
	Host: localhost:8080
	Connection: keep-alive
	Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
	User-Agent: Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Ch
	Accept-Encoding: gzip,deflate,sdch
	Accept-Language: en-US,en;q=0.8,de-DE;q=0.6,de;q=0.4

	2013-11-30 14:36:50.105 <-1 506 HTTP/1.1 302 Found
	Location: http://www.google.com/
	Cache-Control: private
	Content-Type: text/html; charset=UTF-8
	Content-Length: 219
	...

Example 3: Log binary communication
-----------------------------------
lupo tries to detect if binary or textual data is transferred and prints the log accordingly.

This is how a sample Websocket exchange looks like (which is mixed plain text / binary):

	2013-11-30 14:36:50.069 ->1 360	GET /cometd HTTP/1.1
	Host: localhost:8081
	Upgrade: websocket
	Connection: Upgrade
	Sec-WebSocket-Key: LlBN4AdNLrSl5Q0ckDaNaA==
	Sec-WebSocket-Version: 13
	Pragma: no-cache
	Cache-Control: no-cache

	2013-11-30 14:36:50.105 <-1 506 HTTP/1.1 101 Switching Protocols
	Upgrade: websocket
	Connection: Upgrade
	Sec-WebSocket-Accept: pRsBu2MWn3KruE+/O0JDGNcnKTk=

	2013-11-30 14:36:50.069 ->1 360	
	00000000  81 fd 4b bf 74 8b 10 c4  56 e2 2f 9d 4e a9 7a 9d  |..K.t...V./.N.z.|
	00000010  58 a9 38 ca 04 fb 24 cd  00 ee 2f fc 1b e5 25 da  |X.8...$.../...%.|
	00000020  17 ff 22 d0 1a df 32 cf  11 f8 69 85 2f a9 3c da  |.."...2...i./.<.|
	00000030  16 f8 24 dc 1f ee 3f 9d  29 a7 69 dc 1c ea 25 d1  |..$...?.).i...%.|
	00000040  11 e7 69 85 56 a4 26 da  00 ea 64 d7 15 e5 2f cc  |..i.V.&...d.../.|
	00000050  1c ea 20 da 56 a7 69 da  0c ff 65 d1 1b ef 2e f6  |.. .V.i...e.....|
	00000060  10 a9 71 9d 10 ee 2a d3  11 f9 66 8f 44 bb 7b 8f  |..q...*...f.D.{.|
	00000070  56 a7 69 c9 11 f9 38 d6  1b e5 69 85 56 ba 65 8f  |V.i...8...i.V.e.|
	00000080  56 f6 16                                          |V..|

	2013-11-30 14:36:50.105 <-1 506 
	00000000  81 7e 00 b3                                       |.~..|
	[{"id":"1","minimumVersion":"1.0","supportedConnectionTypes":["websocket"],"successful":true,"channel":"/meta/handshake","clientId":"61p2mdhr4fws221tetjlzxega8c","version":"1.0"}]

Example 4: Log statistics
-------------------------
lupo prints statistics in a csv format, if you use the `-style stats` flag.

	lupo -from :8080 -to google.com:80 -style stats

	Date;ConnCount;ConnOpened;ConnClosed;Sent;Received;TotalTransferred
	2013-11-30 14:53:30.042;0;0;0;0;0;0
	2013-11-30 14:53:31.042;0;0;0;0;0;0
	2013-11-30 14:53:32.042;2;2;0;497;506;1003
	2013-11-30 14:53:34.042;2;0;0;523;506;1029
	2013-11-30 14:53:35.042;2;0;0;0;0;0
