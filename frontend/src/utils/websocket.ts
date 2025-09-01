import { SystemStats, PacketInfo, WebSocketMessage } from '@/types/api';

interface WebSocketHandlers {
  onStats: (stats: SystemStats) => void;
  onPacket?: (packet: PacketInfo) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
}

let socket: WebSocket | null = null;

export const connectWebSocket = (handlers: WebSocketHandlers): WebSocket | null => {
  try {
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080';
    socket = new WebSocket(`${wsUrl}/api/v1/ws`);

    socket.onopen = () => {
      console.log('WebSocket connected');
      handlers.onConnect?.();
    };

    socket.onclose = () => {
      console.log('WebSocket disconnected');
      handlers.onDisconnect?.();
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      handlers.onError?.(error);
    };

    socket.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        
        switch (message.type) {
          case 'stats':
            handlers.onStats(message.data as SystemStats);
            break;
          case 'packet':
            handlers.onPacket?.(message.data as PacketInfo);
            break;
          default:
            console.warn('Unknown message type:', message.type);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    return socket;
  } catch (error) {
    console.error('Error connecting to WebSocket:', error);
    return null;
  }
};

export const disconnectWebSocket = () => {
  if (socket) {
    socket.close();
    socket = null;
  }
};

export const getWebSocketState = (): number | null => {
  return socket?.readyState ?? null;
};
