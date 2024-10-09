import React from 'react';

import { toZonedTime } from 'date-fns-tz';

import { DateTimeUtils } from '@utils/date.ts';

export const DateCell = ({ value }: { value: string }) => {
  if (!value) {
    return <p className='text-gray-400'>Unknown</p>;
  }
  const date = toZonedTime(value, 'UTC').toUTCString();

  return (
    <p className='text-gray-700 cursor-default truncate'>
      {DateTimeUtils.format(date, DateTimeUtils.dateWithAbreviatedMonth)}
    </p>
  );
};
