"use client";

import { SystemStats } from '@/types/api';

interface StatsCardsProps {
  stats: SystemStats | null;
}

export default function StatsCards({ stats }: StatsCardsProps) {
  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  };

  const formatUptime = (nanoseconds: number): string => {
    const seconds = Math.floor(nanoseconds / 1000000000);
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const cards = [
    {
      title: 'Total Packets',
      value: stats ? stats.totalPackets.toLocaleString() : '0',
      icon: '📊',
      color: 'bg-blue-50 text-blue-600 border-blue-200',
    },
    {
      title: 'Total Bytes',
      value: stats ? formatBytes(stats.totalBytes) : '0 B',
      icon: '💾',
      color: 'bg-green-50 text-green-600 border-green-200',
    },
    {
      title: 'Active Interfaces',
      value: stats ? stats.activeInterfaces.toString() : '0',
      icon: '🌐',
      color: 'bg-purple-50 text-purple-600 border-purple-200',
    },
    {
      title: 'Uptime',
      value: stats ? formatUptime(stats.uptime) : '00:00:00',
      icon: '⏱️',
      color: 'bg-orange-50 text-orange-600 border-orange-200',
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {cards.map((card, index) => (
        <div
          key={index}
          className={`p-6 rounded-lg border-2 ${card.color} transition-all hover:shadow-md`}
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium opacity-75">{card.title}</p>
              <p className="text-2xl font-bold mt-1">{card.value}</p>
            </div>
            <div className="text-3xl opacity-75">{card.icon}</div>
          </div>
        </div>
      ))}
    </div>
  );
}
