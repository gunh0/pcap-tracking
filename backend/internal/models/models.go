package models

import (
	"time"
)

// PacketStats represents packet statistics for an interface
type PacketStats struct {
	Interface     string            `json:"interface"`
	Timestamp     time.Time         `json:"timestamp"`
	PacketsTotal  uint64            `json:"packetsTotal"`
	BytesTotal    uint64            `json:"bytesTotal"`
	PacketsPerSec uint64            `json:"packetsPerSec"`
	BytesPerSec   uint64            `json:"bytesPerSec"`
	ProtocolStats map[string]uint64 `json:"protocolStats"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Addresses   []string `json:"addresses"`
	IsUp        bool     `json:"isUp"`
	IsLoopback  bool     `json:"isLoopback"`
}

// PacketInfo represents detailed packet information
type PacketInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Interface string    `json:"interface"`
	Length    int       `json:"length"`
	Protocol  string    `json:"protocol"`
	SrcIP     string    `json:"srcIp"`
	DstIP     string    `json:"dstIp"`
	SrcPort   int       `json:"srcPort,omitempty"`
	DstPort   int       `json:"dstPort,omitempty"`
	Info      string    `json:"info,omitempty"`
}

// InterfaceStatus represents the current status of an interface
type InterfaceStatus struct {
	Name      string     `json:"name"`
	IsActive  bool       `json:"isActive"`
	Error     string     `json:"error,omitempty"`
	StartTime *time.Time `json:"startTime,omitempty"`
}

// SystemStats represents overall system statistics
type SystemStats struct {
	Timestamp        time.Time               `json:"timestamp"`
	TotalPackets     uint64                  `json:"totalPackets"`
	TotalBytes       uint64                  `json:"totalBytes"`
	ActiveInterfaces int                     `json:"activeInterfaces"`
	InterfaceStats   map[string]*PacketStats `json:"interfaceStats"`
	Uptime           int64                   `json:"uptime"` // Duration in nanoseconds
}

// CaptureFilter represents a packet capture filter
type CaptureFilter struct {
	Interface string `json:"interface,omitempty"`
	Protocol  string `json:"protocol,omitempty"`
	SrcIP     string `json:"srcIp,omitempty"`
	DstIP     string `json:"dstIp,omitempty"`
	Port      int    `json:"port,omitempty"`
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
