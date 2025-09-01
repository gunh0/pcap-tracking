package config

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// AutoDetectActiveInterface detects the most active network interface
func AutoDetectActiveInterface() (string, error) {
	interfaces, err := getInterfaceStats()
	if err != nil {
		return "", err
	}

	var bestInterface string
	var maxPackets int64

	for name, packets := range interfaces {
		// Skip loopback and virtual interfaces
		if strings.HasPrefix(name, "lo") ||
			strings.HasPrefix(name, "utun") ||
			strings.HasPrefix(name, "gif") ||
			strings.HasPrefix(name, "stf") ||
			strings.HasPrefix(name, "XHC") ||
			strings.HasPrefix(name, "awdl") {
			continue
		}

		// Check if interface is up
		if iface, err := net.InterfaceByName(name); err == nil {
			if iface.Flags&net.FlagUp != 0 && packets > maxPackets {
				maxPackets = packets
				bestInterface = name
			}
		}
	}

	if bestInterface == "" {
		return "en0", nil // Fallback to en0
	}

	return bestInterface, nil
}

// getInterfaceStats gets network interface statistics using netstat
func getInterfaceStats() (map[string]int64, error) {
	cmd := exec.Command("netstat", "-i")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat: %v", err)
	}

	stats := make(map[string]int64)
	lines := strings.Split(string(output), "\n")

	// Skip header lines
	for i, line := range lines {
		if i < 2 {
			continue
		}

		fields := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(line), -1)
		if len(fields) >= 5 {
			interfaceName := fields[0]
			// Remove * suffix for down interfaces
			interfaceName = strings.TrimSuffix(interfaceName, "*")

			if packetsStr := fields[4]; packetsStr != "-" {
				if packets, err := strconv.ParseInt(packetsStr, 10, 64); err == nil {
					stats[interfaceName] = packets
				}
			}
		}
	}

	return stats, nil
}
