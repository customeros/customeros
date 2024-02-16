import { useState, useEffect } from 'react';

import { Channel, Presence } from 'phoenix';

type Meta = {
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

export const usePresence = (channel: Channel | null) => {
  const [presence, setPresence] = useState<PresenceDiff | null>(null);
  const presentUsers = !presence
    ? []
    : Object.values(presence).flatMap((p) =>
        p?.metas.flatMap((m) => m?.username),
      );

  useEffect(() => {
    if (!channel) return;

    channel.on('presence_state', (response) => {
      const nextPresences = Presence.syncState(
        presence ?? {},
        response,
      ) as PresenceDiff;

      setPresence(nextPresences);
    });

    channel.on('presence_diff', (response) => {
      const nextPresences = Presence.syncDiff(
        presence ?? {},
        response,
      ) as PresenceDiff;

      setPresence(nextPresences);
    });

    return () => {
      channel.off('presence_state');
      channel.off('presence_diff');
    };
  }, [channel]);

  return { presence, presentUsers };
};
