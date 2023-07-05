'use client';
import React, { FC, PropsWithChildren, useState } from 'react';
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
  const [isOpen, setOpen] = useState(false);
  return (
    <Box mt={showDate ? 3 : 4} mr={6}>
      {showDate && (
        <Text color='#667085' fontSize='12px' fontWeight={500} marginBottom={4}>
          {DateTimeUtils.format(date, DateTimeUtils.defaultFormatShortString)}
        </Text>
      )}
      <div onClick={() => setOpen(!isOpen)}>{children}</div>
    </Box>
  );
};
