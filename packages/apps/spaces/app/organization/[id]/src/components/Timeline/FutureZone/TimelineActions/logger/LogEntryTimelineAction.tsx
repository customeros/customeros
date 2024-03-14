import React, { useEffect } from 'react';

import { Box } from '@ui/layout/Box';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';

import { Logger } from './components/Logger';

export const LogEntryTimelineAction: React.FC = () => {
  const { virtuosoRef } = useTimelineRefContext();

  useEffect(() => {
    virtuosoRef?.current?.scrollBy({ top: 300 });
  }, [virtuosoRef]);

  return (
    <Box
      borderRadius={'md'}
      boxShadow={'lg'}
      m={6}
      mt={2}
      p={6}
      pt={4}
      bg='white'
      border='1px solid'
      borderColor='gray.100'
    >
      <Logger />
    </Box>
  );
};
