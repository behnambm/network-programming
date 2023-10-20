package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Create a UDP address to listen on port 5000
	serverAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:5000")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// Create a UDP connection for the server
	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error creating UDP server:", err)
		os.Exit(1)
	}
	defer serverConn.Close()

	fmt.Println("UDP server is listening on port 5000...")

	buffer := make([]byte, 1024)

	for {
		n, addr, err := serverConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			continue
		}

		clientData := buffer[:n]

		// Create a UDP address for the DNS server (1.1.1.1:53)
		dnsAddr, err := net.ResolveUDPAddr("udp", "1.1.1.1:53")
		if err != nil {
			fmt.Println("Error resolving DNS server address:", err)
			continue
		}

		// Create a UDP connection to the DNS server
		dnsConn, err := net.DialUDP("udp", nil, dnsAddr)
		if err != nil {
			fmt.Println("Error connecting to DNS server:", err)
			continue
		}
		defer dnsConn.Close()

		// Forward the data to the DNS server
		_, err = dnsConn.Write(clientData)
		if err != nil {
			fmt.Println("Error sending data to DNS server:", err)
			continue
		}

		// Receive the response from the DNS server
		n, err = dnsConn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving DNS response:", err)
			continue
		}

		// Send the DNS response back to the client
		_, err = serverConn.WriteToUDP(buffer[:n], addr)
		if err != nil {
			fmt.Println("Error sending DNS response to client:", err)
		}
	}
}
