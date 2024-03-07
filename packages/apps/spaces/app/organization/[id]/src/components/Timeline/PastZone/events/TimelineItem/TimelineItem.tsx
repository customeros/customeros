'use client';
import React, { FC, PropsWithChildren } from 'react';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';

interface TimelineItemProps extends PropsWithChildren {
  date: string;
  showDate: boolean;
}

export const TimelineItem: FC<TimelineItemProps> = ({
  date,
  showDate,
  children,
}) => {
  return (
    <Box px={6} pb={showDate ? 2 : 2} bg='gray.25'>
      {showDate && (
        <Text
          color='gray.500'
          fontSize='12px'
          fontWeight={500}
          marginBottom={2}
        >
          {DateTimeUtils.format(date, DateTimeUtils.defaultFormatShortString)}
        </Text>
      )}
      {children}
    </Box>
  );
};
