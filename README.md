lupo. simple network logging.
=============================

lupo is useful for debugging network communications.

lupo is tiny & easy to use. Use it when you don't need or want a full blown wireshark.

lupo listens to a tcp port, and forwards it to another.
All communication which goes through, both ways, is logged to stdout in a convenient format.

lupo just 150 lines of Go code. It runs on all platforms supported by Go. i.e. Linux, Windows, OSX, FreeBSD and more.

Example 1: Log http traffic
---------------------------

	lupo -from :8080 -to google.com:80

If you point your browser to http://localhost:8080, lupo will print something like:

	22:31:01.382 Listening to [:8080], forwarding to [google.com:80]
	22:31:14.818 New Conn: 1
	22:31:14.890 ->1
	GET / HTTP/1.1
	Host: localhost:8080
	Connection: keep-alive
	Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
	User-Agent: Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Ch
	Accept-Encoding: gzip,deflate,sdch
	Accept-Language: en-US,en;q=0.8,de-DE;q=0.6,de;q=0.4

	22:31:14.929 <-1
	HTTP/1.1 302 Found
	Location: http://www.google.com/
	Cache-Control: private
	Content-Type: text/html; charset=UTF-8
	Content-Length: 219
	...


Example 2: Log binary communication
-----------------------------------
lupo tries to detect if binary or textual data is transferred and prints the log accordingly.

This is how a sample Websocket exchange looks like (which is mixed plain text / binary):

	22:46:55.159 ->1
	GET /cometd HTTP/1.1
	Host: localhost:8081
	Upgrade: websocket
	Connection: Upgrade
	Sec-WebSocket-Key: LlBN4AdNLrSl5Q0ckDaNaA==
	Sec-WebSocket-Version: 13
	Pragma: no-cache
	Cache-Control: no-cache

	22:46:55.160 <-1
	HTTP/1.1 101 Switching Protocols
	Upgrade: websocket
	Connection: Upgrade
	Sec-WebSocket-Accept: pRsBu2MWn3KruE+/O0JDGNcnKTk=

	22:46:55.209 ->1
	00000000  81 fd 4b bf 74 8b 10 c4  56 e2 2f 9d 4e a9 7a 9d  |..K.t...V./.N.z.|
	00000010  58 a9 38 ca 04 fb 24 cd  00 ee 2f fc 1b e5 25 da  |X.8...$.../...%.|
	00000020  17 ff 22 d0 1a df 32 cf  11 f8 69 85 2f a9 3c da  |.."...2...i./.<.|
	00000030  16 f8 24 dc 1f ee 3f 9d  29 a7 69 dc 1c ea 25 d1  |..$...?.).i...%.|
	00000040  11 e7 69 85 56 a4 26 da  00 ea 64 d7 15 e5 2f cc  |..i.V.&...d.../.|
	00000050  1c ea 20 da 56 a7 69 da  0c ff 65 d1 1b ef 2e f6  |.. .V.i...e.....|
	00000060  10 a9 71 9d 10 ee 2a d3  11 f9 66 8f 44 bb 7b 8f  |..q...*...f.D.{.|
	00000070  56 a7 69 c9 11 f9 38 d6  1b e5 69 85 56 ba 65 8f  |V.i...8...i.V.e.|
	00000080  56 f6 16                                          |V..|

	22:46:55.212 <-1
	00000000  81 7e 00 b3                                       |.~..|
	[{"id":"1","minimumVersion":"1.0","supportedConnectionTypes":["websocket"],"successful":true,"channel":"/meta/handshake","clientId":"61p2mdhr4fws221tetjlzxega8c","version":"1.0"}]
