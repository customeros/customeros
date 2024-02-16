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
  const socketPath = `${env.REALTIME_WS_PATH}/socket`;

  useEffect(() => {
    const socket = new Socket(socketPath, { params: { token: '123' } });
    socket.connect();
    setSocket(socket);
  }, [socketPath]);

  return (
    <PhoenixSocketContext.Provider value={{ socket }}>
      {children}
    </PhoenixSocketContext.Provider>
  );
};

export { PhoenixSocketContext, PhoenixSocketProvider };
