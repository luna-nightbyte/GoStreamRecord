package utils

import (
	"fmt"
	"log"
	"net"
	"os/user"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Request struct {
	Address string
	Port    string
	Count   int
}

// IsTCPPortInUse checks if a given TCP port on localhost is currently in use.
func (rq Request) IsTCPPortInUse(port int) bool {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// If there's an error, it usually means the address is already in use.
		// You can inspect the error type for more specific handling if needed.
		// For example, checking if it's a net.OpError and its Err field.
		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" {
			// A common error message for "address already in use"
			// might contain "bind: address already in use" or similar.
			// However, relying on the string message is less robust than checking Op.
			fmt.Printf("Port %d is likely in use: %v\n", port, err)
			return true
		}
		// Other errors might indicate a different problem (e.g., permission denied)
		fmt.Printf("Error checking port %d: %v\n", port, err)
		return true // Treat other errors as "in use" to be safe
	}
	defer listener.Close() // Close the listener immediately after successful bind
	fmt.Printf("Port %d is available.\n", port)
	return false
}

// IsUDPPortInUse checks if a given UDP port on localhost is currently in use.
func (rq Request) IsUDPPortInUse(port int) bool {
	address := ":" + strconv.Itoa(port)
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" {
			fmt.Printf("UDP Port %d is likely in use: %v\n", port, err)
			return true
		}
		fmt.Printf("Error checking UDP port %d: %v\n", port, err)
		return true // Treat other errors as "in use" to be safe
	}
	defer conn.Close()
	fmt.Printf("UDP Port %d is available.\n", port)
	return false
}

func Ping(ipAddr string) bool {
	currentUser, err := user.Current()
	if currentUser.Uid != "0" {
		log.Println("Cannot ping without root privileges!")
		return false
	}
	// Resolve the IP address.
	addr, err := net.ResolveIPAddr("ip", ipAddr)
	if err != nil {
		log.Println("Error resolving IP address:", err)
		return false
	}

	// Create a new ICMP connection.
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Println("Error creating ICMP connection:", err)
		return false
	}
	defer conn.Close()

	// Create the ICMP message.
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1, // Can be any value, but should be unique.
			Seq:  1, // Sequence number.
			Data: []byte("ping"),
		},
	}

	// Marshal the message into a byte slice.
	msgBytes, err := message.Marshal(nil)
	if err != nil {
		log.Println("Error marshaling ICMP message:", err)
		return false
	}

	// Set a deadline for the ping request.
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Record the start time.
	// start := time.Now()

	// Send the ICMP message.
	_, err = conn.WriteTo(msgBytes, addr)
	if err != nil {
		log.Println("Error sending ICMP message:", err)
		return false
	}

	// Create a buffer to read the response.
	reply := make([]byte, 1500)

	// Wait for the response.
	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		return false
	}

	// Calculate the round-trip time.
	// rtt := time.Since(start)

	// Unmarshal the response.
	replyMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		log.Println("Error parsing ICMP reply:", err)
		return false
	}

	// Check if the reply is an Echo Reply.
	if replyMsg.Type == ipv4.ICMPTypeEchoReply {
		return true
	} else {
		return false
	}
}

// returns local non-loopback IPv4 address
func GetLocalIp() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagPointToPoint != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				// We only care about IPv4 addresses
				ip := ipnet.IP.To4()
				if ip == nil {
					continue
				}
				if ip.IsLoopback() {
					continue
				}
				if ip.IsPrivate() {
					return ip.String()
				}
			}
		}
	}
	return ""
}
