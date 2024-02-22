'use client';

import { useState, useEffect, createContext } from 'react';

import { Socket } from 'phoenix';

import { useEnv } from '@shared/hooks/useEnv';

const PhoenixSocketContext = createContext<{ socket: Socket | null }>({
  socket: null,
});

const PhoenixSocketProvider = ({ children }: { children: React.ReactNode }) => {
  const env = useEnv();
  const [socket, setSocket] = useState<Socket | null>(null);

  const token = env.REALTIME_WS_API_KEY;
  const socketPath = `${env.REALTIME_WS_PATH}/socket`;

  useEffect(() => {
    if (!token) return;
    try {
      const socket = new Socket(socketPath, {
        params: { token },
      });
      socket.connect();
      setSocket(socket);
    } catch (e) {
      // TODO: log error
    }
  }, [socketPath, token]);

  return (
    <PhoenixSocketContext.Provider value={{ socket }}>
      {children}
    </PhoenixSocketContext.Provider>
  );
};

export { PhoenixSocketContext, PhoenixSocketProvider };
