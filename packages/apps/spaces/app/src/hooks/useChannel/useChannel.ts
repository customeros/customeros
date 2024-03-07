import type { Channel } from 'phoenix';

import { useState, useEffect, useContext } from 'react';

import { Presence } from 'phoenix';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { PhoenixSocketContext } from '../../components/Providers/SocketProvider';

type Meta = {
  color: string;
  phx_ref: string;
  username: string;
  online_at: number;
  metadata: { source: string };
};
type PresenceDiff = {
  [key: string]: {
    metas: Meta[];
  };
};

type PresenceState = { metas: Meta[] }[];

export const useChannel = (channelName: string) => {
  const client = getGraphQLClient();
  const { socket } = useContext(PhoenixSocketContext);

  const [presenceState, setPresenceState] = useState<PresenceState | null>(
    null,
  );

  const [channel, setChannel] = useState<Channel | null>(null);
  const [presence, setPresence] = useState<PresenceDiff | null>(null);
  const presentUsers = parsePresentUsers(presenceState || []);

  const { data } = useGlobalCacheQuery(client);

  const user = data?.global_Cache?.user;
  const user_id = user?.id;
  const username = (() => {
    if (!user) return;
    const fullName = [user?.firstName, user?.lastName].join(' ').trim();
    const email = user?.emails?.[0]?.email;

    return fullName || email;
  })();

  useEffect(() => {
    if (!socket || !user_id) return;

    const phoenixChannel = socket?.channel(channelName, {
      user_id,
      username,
    });

    if (!phoenixChannel) return;

    phoenixChannel
      ?.join()
      ?.receive('ok', () => {
        setChannel(phoenixChannel);
      })
      .receive('error', () => {
        // TODO: handle error
      });

    const presence = new Presence(phoenixChannel);

    presence.onSync(() => {
      setPresenceState(presence.list());
    });

    return () => {
      phoenixChannel.leave();
    };
  }, [setPresence, socket, user_id]);

  return { username, channel, presence, presentUsers };
};

function parsePresentUsers(presenceState: PresenceState) {
  return presenceState.map((p) => [p.metas?.[0].username, p.metas?.[0].color]);
}
