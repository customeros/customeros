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

  constructor() {
    this.http = createHttpClient({});
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
  ): Promise<void | { channel: Channel; latest: LatestDiff | null }> {
    return new Promise((resolve, reject) => {
      const existingChannel = this.channels.get(`${channelName}:${id}`);
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
          this.channels.set(id, channel);
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
