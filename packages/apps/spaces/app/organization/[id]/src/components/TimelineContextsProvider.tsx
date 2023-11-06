'use client';
import React, { PropsWithChildren } from 'react';

import { TimelineRefContextProvider } from '@organization/src/components/Timeline/context/TimelineRefContext';
import { TimelineEventPreviewContextContextProvider } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

export const TimelineContextsProvider: React.FC<
  PropsWithChildren & { id: string }
> = ({ children, id }) => {
  return (
    <TimelineRefContextProvider>
      <TimelineEventPreviewContextContextProvider id={id}>
        {children}
      </TimelineEventPreviewContextContextProvider>
    </TimelineRefContextProvider>
  );
};
