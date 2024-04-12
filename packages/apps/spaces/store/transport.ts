import { makeAutoObservable } from 'mobx';
import { Socket, Channel } from 'phoenix';

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
  channel: Channel | null = null;
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
    this.connect();
  }

  connect() {
    this?.socket?.connect();
  }

  join(channelName: string) {
    if (checkRuntime() === 'node') {
      console.info('Node runtime detected: skipping channel join.');

      return;
    }

    const channel = this?.socket?.channel(channelName, {
      ...this.metadata,
    });

    if (!channel) {
      throw new Error('Channel not found');
    }

    channel
      .join()
      .receive('ok', () => {
        this.channel = channel;
      })
      .receive('error', () => {
        throw new Error('Error joining channel');
      });
  }

  setMetadata(metadata: TransportLayerMetadata) {
    this.metadata = metadata;
  }
}

function checkRuntime() {
  return typeof window === 'undefined' ? 'node' : 'browser';
}
