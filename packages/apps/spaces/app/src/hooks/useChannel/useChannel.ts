import type { Channel } from 'phoenix';

import { useState, useEffect, useContext } from 'react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { usePresence } from '../usePresence';
import { PhoenixSocketContext } from '../../components/Providers/SocketProvider';

export const useChannel = (channelName: string) => {
  const client = getGraphQLClient();
  const { socket } = useContext(PhoenixSocketContext);

  const [channel, setChannel] = useState<Channel | null>(null);
  const { presence, presentUsers } = usePresence(channel);

  const { data, isPending } = useGlobalCacheQuery(client);

  const user = data?.global_Cache?.user;
  const user_id = user?.id;
  const username = (() => {
    if (!user) return;
    const fullName = [user?.firstName, user?.lastName].join(' ').trim();
    const email = user?.emails?.[0]?.email;

    return fullName || email;
  })();

  useEffect(() => {
    if (isPending || !user_id || !username) return;

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

    return () => {
      phoenixChannel.leave();
    };
  }, [channelName, isPending, user_id, username]);

  return { channel, presence, presentUsers };
};
