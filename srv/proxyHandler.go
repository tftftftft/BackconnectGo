package main

import (
	"io"
	"log"
	"net"
)

func Forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	io.Copy(src, dst)
}

func HandleClient(clientConn net.Conn, socks5Addr string) {
	socks5Conn, err := net.Dial("tcp", socks5Addr)
	if err != nil {
		log.Printf("Failed to connect to SOCKS5 proxy: %v", err)
		return
	}

	go Forward(clientConn, socks5Conn)
	go Forward(socks5Conn, clientConn)
}

func StartServer(listenPort string, socks5Addr string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+listenPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", listenPort, err)
	}
	log.Printf("Listening on port %s forwarding to %s", listenPort, socks5Addr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", clientConn.RemoteAddr())
		go HandleClient(clientConn, socks5Addr)
	}
}
