import { makeAutoObservable } from 'mobx';
import { Socket, Channel } from 'phoenix';

import { LatestDiff } from './types';

type TransportLayerMetadata = {
  user_id: string;
  username: string;
};

interface TransportLayerOptions {
  token: string;
  socketPath: string;
}

export class TransportLayer {
  socket: Socket | null = null;
  channels: Map<string, Channel> = new Map();
  metadata: TransportLayerMetadata = {
    user_id: '',
    username: '',
  };

  constructor({ token, socketPath }: TransportLayerOptions) {
    if (checkRuntime() === 'node') {
      console.info(
        'Node runtime detected: skipping TransportLayer initialization.',
      );

      return;
    }

    makeAutoObservable(this);
    this.socket = new Socket(socketPath, {
      params: { token },
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
      if (checkRuntime() === 'node') {
        console.info('Node runtime detected: skipping channel join.');
        resolve();

        return;
      }

      const existingChannel = this.channels.get(id);
      if (existingChannel) {
        resolve({ channel: existingChannel, latest: null });

        return;
      }

      const channel = this?.socket?.channel(`${channelName}:${id}`, {
        ...this.metadata,
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

  setMetadata(metadata: TransportLayerMetadata) {
    this.metadata = metadata;
  }
}

function checkRuntime() {
  return typeof window === 'undefined' ? 'node' : 'browser';
}
