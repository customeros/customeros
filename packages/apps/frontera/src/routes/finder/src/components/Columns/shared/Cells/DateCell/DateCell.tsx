import React from 'react';

import { DateTimeUtils } from '@utils/date.ts';

export const DateCell = ({ value }: { value: string }) => {
  if (!value) {
    return <p className='text-gray-400'>Unknown</p>;
  }

  return (
    <p className='text-gray-700 cursor-default truncate'>
      {DateTimeUtils.format(value, DateTimeUtils.dateWithAbreviatedMonth)}
    </p>
  );
};
