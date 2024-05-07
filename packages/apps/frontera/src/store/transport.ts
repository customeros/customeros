import { Socket, Channel } from 'phoenix';
import axios, { AxiosInstance } from 'axios';
import { GraphQLClient } from 'graphql-request';

import { LatestDiff } from './types';

export interface TransportLayerOptions {
  email: string;
  userId: string;
  sessionToken: string;
}

export class TransportLayer {
  http: AxiosInstance;
  client: GraphQLClient;
  socket: Socket | null = null;
  channels: Map<string, Channel> = new Map();
  channelMeta: Record<string, unknown> = {};
  isAuthenthicated = false;

  constructor(options?: TransportLayerOptions) {
    const headers = {
      Authorization: `Bearer ${options?.sessionToken}`,
      'X-Openline-USERNAME': options?.email ?? '',
    };

    this.client = createGraphqlClient(options ? headers : {});
    this.http = createHttpClient(options ? headers : {});

    if (!options) return;

    this.isAuthenthicated = true;
    this.channelMeta = {
      user_id: options?.userId,
      username: options?.email,
    };

    this.socket = new Socket(import.meta.env.VITE_REALTIME_WS_API_URL ?? '', {
      params: { token: import.meta.env.VITE_REALTIME_WS_API_KEY },
    });

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
      const existingChannel = this.channels.get(id);
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
