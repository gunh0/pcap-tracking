"use client";

import { NetworkInterface, SystemStats } from '@/types/api';

interface InterfaceListProps {
  interfaces: NetworkInterface[];
  stats: SystemStats | null;
}

export default function InterfaceList({ interfaces, stats }: InterfaceListProps) {
  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Network Interfaces</h3>
      
      <div className="space-y-4">
        {interfaces.length > 0 ? (
          interfaces.map((interface_, index) => {
            const interfaceStats = stats?.interfaceStats[interface_.name];
            const isActive = interfaceStats !== undefined;
            
            return (
              <div
                key={index}
                className={`p-4 rounded-lg border-2 transition-all ${
                  isActive 
                    ? 'border-green-200 bg-green-50' 
                    : interface_.isUp 
                    ? 'border-gray-200 bg-gray-50' 
                    : 'border-red-200 bg-red-50'
                }`}
              >
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center space-x-3">
                    <div className={`w-3 h-3 rounded-full ${
                      isActive 
                        ? 'bg-green-500 animate-pulse' 
                        : interface_.isUp 
                        ? 'bg-yellow-500' 
                        : 'bg-red-500'
                    }`}></div>
                    <div>
                      <h4 className="font-medium text-gray-900">{interface_.name}</h4>
                      <p className="text-xs text-gray-500">{interface_.description || 'No description'}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                      isActive
                        ? 'bg-green-100 text-green-800'
                        : interface_.isUp
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {isActive ? 'Monitoring' : interface_.isUp ? 'Ready' : 'Down'}
                    </span>
                  </div>
                </div>

                {/* Interface addresses */}
                {interface_.addresses.length > 0 && (
                  <div className="mb-3">
                    <div className="flex flex-wrap gap-1">
                      {interface_.addresses.map((address, addrIndex) => (
                        <span
                          key={addrIndex}
                          className="px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded"
                        >
                          {address}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Statistics */}
                {interfaceStats && (
                  <div className="grid grid-cols-2 gap-4 mt-3 pt-3 border-t border-gray-200">
                    <div>
                      <p className="text-xs text-gray-500">Packets</p>
                      <p className="text-sm font-medium text-gray-900">
                        {interfaceStats.packetsTotal.toLocaleString()}
                      </p>
                      <p className="text-xs text-green-600">
                        {interfaceStats.packetsPerSec}/sec
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-500">Data</p>
                      <p className="text-sm font-medium text-gray-900">
                        {formatBytes(interfaceStats.bytesTotal)}
                      </p>
                      <p className="text-xs text-green-600">
                        {formatBytes(interfaceStats.bytesPerSec)}/sec
                      </p>
                    </div>
                  </div>
                )}

                {/* Interface properties */}
                <div className="flex space-x-4 mt-2">
                  {interface_.isLoopback && (
                    <span className="text-xs text-blue-600">Loopback</span>
                  )}
                  {!interface_.isUp && (
                    <span className="text-xs text-red-600">Interface Down</span>
                  )}
                </div>
              </div>
            );
          })
        ) : (
          <div className="text-center py-8 text-gray-500">
            <p>No network interfaces found</p>
            <p className="text-sm mt-1">Check if the backend service is running</p>
          </div>
        )}
      </div>
    </div>
  );
}
