import { Socket, Channel } from 'phoenix';
import axios, { AxiosInstance } from 'axios';
import { GraphQLClient } from 'graphql-request';

import { LatestDiff } from './types';

export interface TransportOptions {
  email: string;
  userId: string;
  sessionToken: string;
}

export class Transport {
  http: AxiosInstance;
  graphql: GraphQLClient;
  socket: Socket | null = null;
  refId: string = crypto.randomUUID();
  channels: Map<string, Channel> = new Map();
  channelMeta: Record<string, unknown> = {};
  stream: ReturnType<typeof createStreamClient>;

  constructor() {
    this.http = createHttpClient({});
    this.stream = createStreamClient({});
    this.graphql = createGraphqlClient({});

    this.socket = new Socket(
      `${import.meta.env.VITE_REALTIME_WS_PATH}/socket`,
      {
        params: { token: import.meta.env.VITE_REALTIME_WS_API_KEY },
      },
    );

    if (this.socket.isConnected()) return;
    this.connect();
  }

  connect() {
    this?.socket?.connect();
  }

  join(
    channelName: string,
    id: string,
    version: number,
    group?: boolean,
  ): Promise<void | { channel: Channel; latest: LatestDiff | null }> {
    return new Promise((resolve, reject) => {
      const channelKey = group ? channelName : id;
      const existingChannel = this.channels.get(channelKey);

      if (existingChannel) {
        resolve({ channel: existingChannel, latest: null });

        return;
      }

      const channel = this?.socket?.channel(`${channelName}:${id}`, {
        ...this.channelMeta,
        version,
      });

      if (!channel) {
        reject(new Error('Channel not found'));

        return;
      }

      channel
        .join()
        .receive('ok', (res: LatestDiff) => {
          this.channels.set(channelKey, channel);
          resolve({ latest: res, channel });
        })
        .receive('error', () => {
          reject(new Error('Error joining channel'));
        });
    });
  }

  leaveChannel(id: string) {
    const channel = this.channels.get(id);

    if (!channel) {
      return;
    }

    channel.leave();
    this.channels.delete(id);
  }

  disconnect() {
    this.channels.forEach((channel) => {
      channel.leave();
    });
    this?.socket?.disconnect();
  }

  setChannelMeta(meta: Record<string, unknown>) {
    this.channelMeta = meta;
  }

  setHeaders(headers: Record<string, string>) {
    this.http = createHttpClient(headers);
    this.stream = createStreamClient(headers);
    this.graphql = createGraphqlClient(headers);
  }
}

function createHttpClient(headers?: Record<string, string>) {
  const instance = axios.create({
    baseURL: import.meta.env.VITE_MIDDLEWARE_API_URL,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  });

  return instance;
}

function createGraphqlClient(headers?: Record<string, string>) {
  return new GraphQLClient(
    `${import.meta.env.VITE_MIDDLEWARE_API_URL}/customer-os-api`,
    {
      headers,
    },
  );
}

function createStreamClient(headers?: Record<string, string>) {
  const baseUrl = import.meta.env.VITE_MIDDLEWARE_API_URL;

  return async <TData extends object>(
    endpoint: string,
    options: RequestInit & { onData?: (data: TData) => void },
  ) => {
    const response = await fetch(`${baseUrl}/customer-os-stream${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...headers,
        ...options.headers,
      },
    });

    if (!response.ok) {
      // Handle HTTP errors
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    let result;
    let buffer = '';

    while (!(result = await reader?.read())?.done) {
      buffer += decoder.decode(result?.value, { stream: true });

      let boundary = buffer.indexOf('\n');

      while (boundary !== -1) {
        const completeChunk = buffer.substring(0, boundary);

        buffer = buffer.substring(boundary + 1);
        boundary = buffer.indexOf('\n');

        if (completeChunk) {
          try {
            const data = JSON.parse(completeChunk);

            options.onData?.(data);
          } catch (e) {
            console.error('Error parsing JSON:', e);
          }
        }
      }
    }
  };
}
