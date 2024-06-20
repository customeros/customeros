import React, { useEffect } from 'react';

import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';

import { Logger } from './components/Logger';

export const LogEntryTimelineAction: React.FC = () => {
  const { virtuosoRef } = useTimelineRefContext();

  useEffect(() => {
    virtuosoRef?.current?.scrollBy({ top: 300 });
  }, [virtuosoRef]);

  return (
    <div className='rounded-md shadow-lg m-6 mt-2 p-6 pt-4 bg-white border border-gray-100 max-w-[800px]'>
      <Logger />
    </div>
  );
};
