'use client';
import React, { FC, PropsWithChildren } from 'react';
import { DateTimeUtils } from '@spaces/utils/date';
import { Text } from '@ui/typography/Text';
import { Box } from '@ui/layout/Box';

interface TimelineItemProps extends PropsWithChildren {
  showDate: boolean;
  date: string;
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
