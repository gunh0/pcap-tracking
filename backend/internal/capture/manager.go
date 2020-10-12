package capture

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"pcap-tracking/backend/internal/config"
	"pcap-tracking/backend/internal/models"
	"pcap-tracking/backend/internal/websocket"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
)

// Manager manages packet capture for multiple interfaces
type Manager struct {
	config        config.CaptureConfig
	hub           *websocket.Hub
	captures      map[string]*InterfaceCapture
	stats         map[string]*models.PacketStats
	systemStats   *models.SystemStats
	mutex         sync.RWMutex
	running       bool
	startTime     time.Time
	statsTicker   *time.Ticker
	lastStatsTime time.Time
}

// InterfaceCapture represents packet capture on a single interface
type InterfaceCapture struct {
	name      string
	handle    *pcap.Handle
	source    *gopacket.PacketSource
	running   bool
	startTime time.Time
	mutex     sync.RWMutex
}

// NewManager creates a new capture manager
func NewManager(config config.CaptureConfig, hub *websocket.Hub) *Manager {
	return &Manager{
		config:   config,
		hub:      hub,
		captures: make(map[string]*InterfaceCapture),
		stats:    make(map[string]*models.PacketStats),
		systemStats: &models.SystemStats{
			InterfaceStats: make(map[string]*models.PacketStats),
		},
		startTime:     time.Now().UTC(),
		lastStatsTime: time.Now().UTC(),
	}
}

// Start starts packet capture on all configured interfaces
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.running {
		return fmt.Errorf("capture manager already running")
	}

	// Start capture on each interface
	for _, ifaceName := range m.config.Interfaces {
		if err := m.startInterface(ifaceName); err != nil {
			logrus.WithFields(logrus.Fields{
				"interface": ifaceName,
				"error":     err,
			}).Error("Failed to start capture on interface")
			continue
		}
	}

	m.running = true

	// Start statistics collection
	m.statsTicker = time.NewTicker(time.Duration(m.config.StatsInterval) * time.Second)
	go m.collectStats()

	logrus.WithField("interfaces", len(m.captures)).Info("Packet capture started")
	return nil
}

// Stop stops packet capture on all interfaces
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running {
		return
	}

	m.running = false

	// Stop statistics collection
	if m.statsTicker != nil {
		m.statsTicker.Stop()
	}

	// Stop all interface captures
	for name := range m.captures {
		m.stopInterface(name)
	}

	logrus.Info("Packet capture stopped")
}

// GetInterfaces returns available network interfaces
func (m *Manager) GetInterfaces() ([]*models.NetworkInterface, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	interfaces := make([]*models.NetworkInterface, 0, len(devices))
	for _, device := range devices {
		iface := &models.NetworkInterface{
			Name:        device.Name,
			Description: device.Description,
			Addresses:   make([]string, 0, len(device.Addresses)),
		}

		for _, addr := range device.Addresses {
			iface.Addresses = append(iface.Addresses, addr.IP.String())
		}

		// Check if interface is up
		if netIface, err := net.InterfaceByName(device.Name); err == nil {
			iface.IsUp = netIface.Flags&net.FlagUp != 0
			iface.IsLoopback = netIface.Flags&net.FlagLoopback != 0
		}

		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

// GetStats returns current statistics
func (m *Manager) GetStats() *models.SystemStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create a copy to avoid race conditions
	stats := &models.SystemStats{
		Timestamp:        time.Now().UTC(),
		TotalPackets:     m.systemStats.TotalPackets,
		TotalBytes:       m.systemStats.TotalBytes,
		ActiveInterfaces: len(m.captures),
		InterfaceStats:   make(map[string]*models.PacketStats),
		Uptime:           int64(time.Since(m.startTime)),
	}

	for name, stat := range m.stats {
		statCopy := *stat
		stats.InterfaceStats[name] = &statCopy
	}

	return stats
}

// startInterface starts packet capture on a specific interface
func (m *Manager) startInterface(ifaceName string) error {
	// Open device
	handle, err := pcap.OpenLive(ifaceName, m.config.SnapLength, m.config.Promiscuous,
		time.Duration(m.config.TimeoutMs)*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %v", ifaceName, err)
	}

	// Set buffer size
	if err := handle.SetBPFFilter(""); err != nil {
		handle.Close()
		return fmt.Errorf("failed to set BPF filter: %v", err)
	}

	capture := &InterfaceCapture{
		name:      ifaceName,
		handle:    handle,
		source:    gopacket.NewPacketSource(handle, handle.LinkType()),
		running:   true,
		startTime: time.Now().UTC(),
	}

	m.captures[ifaceName] = capture
	m.stats[ifaceName] = &models.PacketStats{
		Interface:     ifaceName,
		Timestamp:     time.Now().UTC(),
		ProtocolStats: make(map[string]uint64),
	}

	// Start packet processing goroutine
	go m.processPackets(capture)

	logrus.WithField("interface", ifaceName).Info("Started packet capture")
	return nil
}

// stopInterface stops packet capture on a specific interface
func (m *Manager) stopInterface(ifaceName string) {
	capture, exists := m.captures[ifaceName]
	if !exists {
		return
	}

	capture.mutex.Lock()
	capture.running = false
	capture.mutex.Unlock()

	capture.handle.Close()
	delete(m.captures, ifaceName)
	delete(m.stats, ifaceName)

	logrus.WithField("interface", ifaceName).Info("Stopped packet capture")
}

// processPackets processes packets from an interface
func (m *Manager) processPackets(capture *InterfaceCapture) {
	for {
		capture.mutex.RLock()
		running := capture.running
		capture.mutex.RUnlock()

		if !running {
			break
		}

		packet, err := capture.source.NextPacket()
		if err != nil {
			// Handle timeout errors gracefully - don't log as errors since they're normal
			if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "Timeout") {
				continue
			}
			// Only log other types of errors
			logrus.WithFields(logrus.Fields{
				"interface": capture.name,
				"error":     err,
			}).Warn("Error reading packet")
			continue
		}

		m.handlePacket(capture.name, packet)
	}
}

// handlePacket processes a single packet
func (m *Manager) handlePacket(ifaceName string, packet gopacket.Packet) {
	m.mutex.Lock()
	stats := m.stats[ifaceName]
	stats.PacketsTotal++
	stats.BytesTotal += uint64(len(packet.Data()))
	m.systemStats.TotalPackets++
	m.systemStats.TotalBytes += uint64(len(packet.Data()))

	// Update protocol statistics
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		protocol := networkLayer.LayerType().String()
		stats.ProtocolStats[protocol]++
	}
	m.mutex.Unlock()

	// Create packet info for detailed tracking
	packetInfo := m.extractPacketInfo(ifaceName, packet)
	if packetInfo != nil {
		// Broadcast packet info via WebSocket (with rate limiting)
		if m.hub != nil {
			go m.hub.BroadcastPacket(packetInfo)
		}
	}
}

// extractPacketInfo extracts relevant information from a packet
func (m *Manager) extractPacketInfo(ifaceName string, packet gopacket.Packet) *models.PacketInfo {
	info := &models.PacketInfo{
		Timestamp: packet.Metadata().Timestamp,
		Interface: ifaceName,
		Length:    len(packet.Data()),
	}

	// Extract network layer information
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		info.Protocol = networkLayer.LayerType().String()

		switch layer := networkLayer.(type) {
		case *layers.IPv4:
			info.SrcIP = layer.SrcIP.String()
			info.DstIP = layer.DstIP.String()
		case *layers.IPv6:
			info.SrcIP = layer.SrcIP.String()
			info.DstIP = layer.DstIP.String()
		}
	}

	// Extract transport layer information
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		switch layer := transportLayer.(type) {
		case *layers.TCP:
			info.SrcPort = int(layer.SrcPort)
			info.DstPort = int(layer.DstPort)
			info.Protocol = "TCP"
		case *layers.UDP:
			info.SrcPort = int(layer.SrcPort)
			info.DstPort = int(layer.DstPort)
			info.Protocol = "UDP"
		}
	}

	return info
}

// collectStats periodically collects and broadcasts statistics
func (m *Manager) collectStats() {
	for range m.statsTicker.C {
		m.mutex.RLock()
		running := m.running
		m.mutex.RUnlock()

		if !running {
			break
		}

		// Calculate rates
		now := time.Now().UTC()
		timeDiff := now.Sub(m.lastStatsTime).Seconds()

		m.mutex.Lock()
		for _, stats := range m.stats {
			if timeDiff > 0 {
				// Calculate packets per second and bytes per second
				stats.PacketsPerSec = uint64(float64(stats.PacketsTotal) / timeDiff)
				stats.BytesPerSec = uint64(float64(stats.BytesTotal) / timeDiff)
				stats.Timestamp = now
			}
		}
		m.lastStatsTime = now
		m.mutex.Unlock()

		stats := m.GetStats()
		if m.hub != nil {
			m.hub.BroadcastStats(stats)
		}

		// Log periodic stats summary (only if there's traffic)
		if stats.TotalPackets > 0 {
			logrus.WithFields(logrus.Fields{
				"total_packets": stats.TotalPackets,
				"total_bytes":   stats.TotalBytes,
				"interfaces":    stats.ActiveInterfaces,
			}).Debug("Packet capture statistics")
		}
	}
}
