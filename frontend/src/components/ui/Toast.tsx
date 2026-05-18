'use client';

import { useEffect, useState, useRef } from 'react';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastMessage {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
  duration?: number;
}

interface ToastProps extends ToastMessage {
  onClose: (id: string) => void;
}

const bgColorMap: Record<ToastType, string> = {
  success: 'bg-green-50',
  error: 'bg-red-50',
  info: 'bg-blue-50',
  warning: 'bg-yellow-50',
};

const borderColorMap: Record<ToastType, string> = {
  success: 'border-green-200',
  error: 'border-red-200',
  info: 'border-blue-200',
  warning: 'border-yellow-200',
};

const titleColorMap: Record<ToastType, string> = {
  success: 'text-green-900',
  error: 'text-red-900',
  info: 'text-blue-900',
  warning: 'text-yellow-900',
};

const textColorMap: Record<ToastType, string> = {
  success: 'text-green-800',
  error: 'text-red-800',
  info: 'text-blue-800',
  warning: 'text-yellow-800',
};

const progressBgMap: Record<ToastType, string> = {
  success: 'bg-green-400',
  error: 'bg-red-400',
  info: 'bg-blue-400',
  warning: 'bg-yellow-400',
};

const iconMap: Record<ToastType, string> = {
  success: '✓',
  error: '✕',
  info: 'ⓘ',
  warning: '⚠',
};

export function Toast({
  id,
  type,
  title,
  message,
  duration = 5000,
  onClose,
}: ToastProps) {
  const [timeLeft, setTimeLeft] = useState(duration);
  const [isHovering, setIsHovering] = useState(false);
  const timerRef = useRef<NodeJS.Timeout | undefined>(undefined);

  useEffect(() => {
    if (isHovering) {
      if (timerRef.current) clearInterval(timerRef.current);
      return;
    }

    if (timeLeft <= 0) {
      onClose(id);
      return;
    }

    timerRef.current = setInterval(() => {
      setTimeLeft((prev) => Math.max(0, prev - 100));
    }, 100);

    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [timeLeft, isHovering, id, onClose]);

  return (
    <div
      onMouseEnter={() => setIsHovering(true)}
      onMouseLeave={() => setIsHovering(false)}
      className={`relative mb-3 overflow-hidden rounded-lg border ${bgColorMap[type]} ${borderColorMap[type]} p-4 shadow-md transition-all duration-200`}
    >
      <div className="flex gap-3">
        <div className={`flex-shrink-0 text-lg font-bold ${titleColorMap[type]}`}>
          {iconMap[type]}
        </div>
        <div className="flex-1 min-w-0">
          {title && (
            <h3 className={`font-semibold ${titleColorMap[type]}`}>{title}</h3>
          )}
          <p className={`text-sm ${textColorMap[type]} break-words`}>{message}</p>
        </div>
        <button
          onClick={() => onClose(id)}
          className={`flex-shrink-0 ml-2 text-lg font-bold ${titleColorMap[type]} hover:opacity-70`}
        >
          ×
        </button>
      </div>

      {/* Progress bar */}
      <div className="absolute bottom-0 left-0 right-0 h-1 bg-gray-200 overflow-hidden">
        <div
          className={`h-full ${progressBgMap[type]} transition-all`}
          style={{
            width: `${(timeLeft / duration) * 100}%`,
            transitionDuration: isHovering ? '0ms' : '100ms',
          }}
        />
      </div>
    </div>
  );
}
