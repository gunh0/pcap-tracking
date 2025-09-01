"use client";

import { useState, useEffect } from 'react';

interface ClockProps {
  showDate?: boolean;
  className?: string;
}

export default function Clock({ showDate = true, className = "" }: ClockProps) {
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('en-US', {
      timeZone: 'UTC',
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      timeZone: 'UTC',
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <div className={`text-right ${className}`}>
      {showDate && (
        <div className="text-xs text-gray-500 mb-1">
          {formatDate(currentTime)} UTC
        </div>
      )}
      <div className="text-sm font-mono font-medium text-gray-700">
        {formatTime(currentTime)} UTC
      </div>
    </div>
  );
}
