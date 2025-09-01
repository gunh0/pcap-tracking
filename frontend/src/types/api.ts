export interface NetworkInterface {
  name: string;
  description: string;
  addresses: string[];
  isUp: boolean;
  isLoopback: boolean;
}

export interface PacketStats {
  interface: string;
  timestamp: string;
  packetsTotal: number;
  bytesTotal: number;
  packetsPerSec: number;
  bytesPerSec: number;
  protocolStats: Record<string, number>;
}

export interface SystemStats {
  timestamp: string;
  totalPackets: number;
  totalBytes: number;
  activeInterfaces: number;
  interfaceStats: Record<string, PacketStats>;
  uptime: number;
}

export interface PacketInfo {
  timestamp: string;
  interface: string;
  length: number;
  protocol: string;
  srcIp: string;
  dstIp: string;
  srcPort?: number;
  dstPort?: number;
  info?: string;
}

export interface WebSocketMessage {
  type: 'stats' | 'packet';
  data: SystemStats | PacketInfo;
}

export interface InterfaceStatus {
  name: string;
  isActive: boolean;
  error?: string;
  startTime?: string;
}
