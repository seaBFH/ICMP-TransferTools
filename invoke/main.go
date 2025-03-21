package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	// Parse command line arguments
	ipAddress := flag.String("ip", "", "IP address to download from")
	outputFile := flag.String("output", "", "Output file path")
	flag.Parse()

	if *ipAddress == "" || *outputFile == "" {
		fmt.Println("Usage: program -ip <IP address> -output <output file>")
		os.Exit(1)
	}

	// Call the ICMP download function
	err := InvokeIcmpDownload(*ipAddress, *outputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// TODO: Implement the InvokeIcmpUpload function

// InvokeIcmpDownload downloads a file using ICMP protocol
func InvokeIcmpDownload(ipAddress string, outputPath string) error {
	// Create a new ICMP connection
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("error creating ICMP connection: %v", err)
	}
	defer conn.Close()

	// Open the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	fmt.Println("Downloading file, please wait...")

	// Create a buffer for receiving data
	buf := make([]byte, 1500)

	// Set timeout for the connection
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	for {
		// Create an empty ICMP echo request message
		m := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  1,
				Data: []byte(""),
			},
		}

		// Marshal the message into binary
		b, err := m.Marshal(nil)
		if err != nil {
			return fmt.Errorf("error marshaling ICMP message: %v", err)
		}

		// Send the ICMP packet
		_, err = conn.WriteTo(b, &net.IPAddr{IP: net.ParseIP(ipAddress)})
		if err != nil {
			return fmt.Errorf("error sending ICMP packet: %v", err)
		}

		// Wait for a response
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return fmt.Errorf("timeout waiting for ICMP response")
			}
			return fmt.Errorf("error reading ICMP response: %v", err)
		}

		// Parse the response
		msg, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buf[:n])
		if err != nil {
			return fmt.Errorf("error parsing ICMP message: %v", err)
		}

		// Check if this is an echo reply
		if msg.Type == ipv4.ICMPTypeEchoReply {
			// Extract the data from the reply
			if echoReply, ok := msg.Body.(*icmp.Echo); ok {
				if len(echoReply.Data) > 0 {
					// Convert the data to a string to check if it's the termination signal
					responseStr := string(echoReply.Data)
					if responseStr == "done" {
						fmt.Println("File transfer complete; EXITING.")
						break
					}

					// Write the data to the output file
					_, err = outputFile.Write(echoReply.Data)
					if err != nil {
						return fmt.Errorf("error writing to output file: %v", err)
					}
				}
			}
		}

		// Reset the deadline for the next iteration
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	}

	return nil
}
