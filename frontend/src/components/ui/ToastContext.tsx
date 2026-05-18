'use client';

import React, { createContext, useContext, useState, useCallback } from 'react';
import { ToastMessage, ToastType } from './Toast';

interface ToastContextType {
  toasts: ToastMessage[];
  addToast: (
    message: string,
    options?: {
      type?: ToastType;
      title?: string;
      duration?: number;
    }
  ) => string;
  removeToast: (id: string) => void;
  clearAll: () => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  const addToast = useCallback(
    (
      message: string,
      options?: {
        type?: ToastType;
        title?: string;
        duration?: number;
      }
    ) => {
      const id = `toast-${Date.now()}-${Math.random()}`;
      const toast: ToastMessage = {
        id,
        type: options?.type || 'info',
        title: options?.title,
        message,
        duration: options?.duration || 5000,
      };

      setToasts((prev) => [...prev, toast]);
      return id;
    },
    []
  );

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const clearAll = useCallback(() => {
    setToasts([]);
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast, clearAll }}>
      {children}
    </ToastContext.Provider>
  );
}

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within ToastProvider');
  }
  return context;
}
