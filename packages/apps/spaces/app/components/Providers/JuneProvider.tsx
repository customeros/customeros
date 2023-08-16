'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';

import { useJune } from '@spaces/hooks/useJune';
import { useSession } from 'next-auth/react';

export const JuneProvider = ({ children }: { children: React.ReactNode }) => {
  const analytics = useJune();
  const pathname = usePathname();

  const { data: session } = useSession();

  useEffect(() => {
    if (analytics && session) {
      analytics.user().then((user) => {
        if (!user || user.id() === null) {
          analytics?.identify(session.user?.email, {
            email: session.user?.email,
          });
        }
      });
      if (pathname) {
        analytics.pageView(pathname);
      }
    }
  }, [session, analytics, pathname]);

  return <>{children}</>;
};
