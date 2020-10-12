"use client";

import { useState, useEffect } from 'react';
import { SystemStats } from '@/types/api';

interface NetworkChartProps {
  stats: SystemStats | null;
  isConnected: boolean;
}

interface DataPoint {
  timestamp: string;
  packets: number;
  bytes: number;
}

export default function NetworkChart({ stats, isConnected }: NetworkChartProps) {
  const [dataPoints, setDataPoints] = useState<DataPoint[]>([]);
  const maxDataPoints = 50;

  useEffect(() => {
    if (stats) {
      const newDataPoint: DataPoint = {
        timestamp: new Date(stats.timestamp).toLocaleTimeString('en-US', {
          timeZone: 'UTC',
          hour12: false,
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit'
        }),
        packets: stats.totalPackets,
        bytes: stats.totalBytes,
      };

      setDataPoints(prev => {
        const updated = [...prev, newDataPoint];
        return updated.slice(-maxDataPoints);
      });
    }
  }, [stats]);

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  };

  const maxPackets = Math.max(...dataPoints.map(d => d.packets), 1);
  const maxBytes = Math.max(...dataPoints.map(d => d.bytes), 1);

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900">Network Traffic</h3>
        <div className={`flex items-center space-x-2 ${isConnected ? 'text-green-600' : 'text-red-600'}`}>
          <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-600 animate-pulse' : 'bg-red-600'}`}></div>
          <span className="text-sm">Live</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Packets Chart */}
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">Packet Count</h4>
          <div className="h-32 flex items-end space-x-1">
            {dataPoints.length > 0 ? (
              dataPoints.slice(-20).map((point, index) => {
                const height = (point.packets / maxPackets) * 100;
                return (
                  <div
                    key={index}
                    className="flex-1 bg-blue-500 min-h-[4px] rounded-t transition-all duration-300 hover:bg-blue-600"
                    style={{ height: `${Math.max(height, 4)}%` }}
                    title={`${point.timestamp}: ${point.packets.toLocaleString()} packets`}
                  ></div>
                );
              })
            ) : (
              <div className="flex-1 flex items-center justify-center h-full text-gray-400">
                No data available
              </div>
            )}
          </div>
          <div className="text-xs text-gray-500 mt-2">
            Current: {stats ? stats.totalPackets.toLocaleString() : '0'} packets
          </div>
        </div>

        {/* Bytes Chart */}
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">Data Volume</h4>
          <div className="h-32 flex items-end space-x-1">
            {dataPoints.length > 0 ? (
              dataPoints.slice(-20).map((point, index) => {
                const height = (point.bytes / maxBytes) * 100;
                return (
                  <div
                    key={index}
                    className="flex-1 bg-green-500 min-h-[4px] rounded-t transition-all duration-300 hover:bg-green-600"
                    style={{ height: `${Math.max(height, 4)}%` }}
                    title={`${point.timestamp}: ${formatBytes(point.bytes)}`}
                  ></div>
                );
              })
            ) : (
              <div className="flex-1 flex items-center justify-center h-full text-gray-400">
                No data available
              </div>
            )}
          </div>
          <div className="text-xs text-gray-500 mt-2">
            Current: {stats ? formatBytes(stats.totalBytes) : '0 B'}
          </div>
        </div>
      </div>

      {/* Timeline */}
      <div className="mt-4 flex justify-between text-xs text-gray-400">
        <span>{dataPoints.length > 0 ? `${dataPoints[0]?.timestamp} UTC` : '--:--:-- UTC'}</span>
        <span>{dataPoints.length > 0 ? `${dataPoints[dataPoints.length - 1]?.timestamp} UTC` : '--:--:-- UTC'}</span>
      </div>
    </div>
  );
}
