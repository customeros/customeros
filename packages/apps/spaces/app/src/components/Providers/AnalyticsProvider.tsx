'use client';

import { useEffect } from 'react';

import { useSession } from 'next-auth/react';
import { H } from '@highlight-run/next/client';

declare const heap: {
  identify: (email: string | null | undefined) => void;
  addUserProperties: (
    properties: Record<string, string | null | undefined>,
  ) => void;
};

export const AnalyticsProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const { data: session } = useSession();

  useEffect(() => {
    if (!session) return;
    if (heap) {
      heap.identify(session.user?.email);
      heap.addUserProperties({
        name: session?.user?.name,
      });
    }

    if (session?.user?.email) {
      H.identify(session?.user?.email, {
        name: session?.user?.name ?? 'Unknown',
        playerIdentityId: session.user.playerIdentityId,
      });
    }
  }, [session]);

  return <>{children}</>;
};
