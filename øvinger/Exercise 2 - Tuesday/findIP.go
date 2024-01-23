package main

import (
	"fmt"
	"net"
)

func main() {
	// Get the list of network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate through each interface
	for _, iface := range interfaces {
		// Only consider interfaces with hardware addresses
		if iface.HardwareAddr != nil {
			// Get the list of unicast addresses for the interface
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// Iterate through each unicast address
			for _, addr := range addrs {
				fmt.Printf("Interface: %s, IP: %s\n", iface.Name, addr)
			}
		}
	}
}
