'use client';

import { useEffect } from 'react';
import { useSession } from 'next-auth/react';
declare const heap: any;

export const AnalyticsProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const { data: session } = useSession();

  useEffect(() => {
    if (heap && session) {
      heap.identify(session.user?.email);
      heap.addUserProperties({
        name: session?.user?.name,
      });
    }
  }, [session]);

  return <>{children}</>;
};
