"use client";

import { useState, useEffect } from 'react';
import Dashboard from '@/components/Dashboard';
import { SystemStats, NetworkInterface } from '@/types/api';
import { connectWebSocket, disconnectWebSocket } from '@/utils/websocket';

export default function Home() {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [interfaces, setInterfaces] = useState<NetworkInterface[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    // Fetch initial data
    fetchInterfaces();
    fetchStats();

    // Connect to WebSocket
    const ws = connectWebSocket({
      onStats: (newStats: SystemStats) => {
        setStats(newStats);
      },
      onConnect: () => setIsConnected(true),
      onDisconnect: () => setIsConnected(false),
    });

    return () => {
      disconnectWebSocket();
    };
  }, []);

  const fetchInterfaces = async () => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/interfaces`);
      const data = await response.json();
      setInterfaces(data.interfaces || []);
    } catch (error) {
      console.error('Failed to fetch interfaces:', error);
    }
  };

  const fetchStats = async () => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/stats`);
      const data = await response.json();
      setStats(data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    }
  };

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            PCAP Tracking Dashboard
          </h1>
          <p className="text-gray-600">
            Real-time network packet monitoring and analysis
          </p>
          <div className="mt-4 flex items-center space-x-4">
            <div className={`flex items-center space-x-2 ${isConnected ? 'text-green-600' : 'text-red-600'}`}>
              <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-600' : 'bg-red-600'}`}></div>
              <span className="text-sm font-medium">
                {isConnected ? 'Connected' : 'Disconnected'}
              </span>
            </div>
          </div>
        </div>

        <Dashboard 
          stats={stats}
          interfaces={interfaces}
          isConnected={isConnected}
        />
      </div>
    </main>
  );
}
