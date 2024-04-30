'use client';
import React, { PropsWithChildren } from 'react';

import { TimelineRefContextProvider } from '@organization/components/Timeline/context/TimelineRefContext';
import { TimelineEventPreviewContextContextProvider } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

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
