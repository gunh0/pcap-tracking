"use client";

import { SystemStats, NetworkInterface } from '@/types/api';
import StatsCards from './StatsCards';
import NetworkChart from './NetworkChart';
import InterfaceList from './InterfaceList';
import Clock from './Clock';

interface DashboardProps {
  stats: SystemStats | null;
  interfaces: NetworkInterface[];
  isConnected: boolean;
}

export default function Dashboard({ stats, interfaces, isConnected }: DashboardProps) {
  return (
    <div className="space-y-6">
      {/* Header with Clock */}
      <div className="flex justify-between items-start">
        <div></div>
        <Clock />
      </div>
      
      {/* Stats Overview Cards */}
      <StatsCards stats={stats} />
      
      {/* Main Grid Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Network Traffic Chart */}
        <div className="lg:col-span-2">
          <NetworkChart stats={stats} isConnected={isConnected} />
        </div>
        
        {/* Interface List */}
        <div className="lg:col-span-1">
          <InterfaceList interfaces={interfaces} stats={stats} />
        </div>
        
        {/* Protocol Statistics */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Protocol Distribution
            </h3>
            {stats && Object.keys(stats.interfaceStats).length > 0 ? (
              <div className="space-y-3">
                {Object.entries(
                  Object.values(stats.interfaceStats).reduce((acc, interfaceStats) => {
                    Object.entries(interfaceStats.protocolStats).forEach(([protocol, count]) => {
                      acc[protocol] = (acc[protocol] || 0) + count;
                    });
                    return acc;
                  }, {} as Record<string, number>)
                ).sort(([, a], [, b]) => b - a).slice(0, 10).map(([protocol, count]) => (
                  <div key={protocol} className="flex justify-between items-center">
                    <span className="text-sm font-medium text-gray-700">{protocol}</span>
                    <span className="text-sm text-gray-500">{count.toLocaleString()}</span>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-500 text-sm">No protocol data available</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
