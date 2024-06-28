import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cake = ({ className, ...props }: IconProps) => (
  <svg
    width='24'
    height='24'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
    xmlns='http://www.w3.org/2000/svg'
  >
    <path
      d='M3.23913 14.9716C5.71739 15.7754 7.36957 14.9716 8.19565 12.9622C9.02174 14.9716 10.2609 15.3735 11.5 15.3735C13.1522 15.3735 13.9783 14.5697 14.8043 12.9622C15.2174 15.3735 18.9348 15.7754 19.7609 14.9716M3.23913 14.9716C2.54535 14.5216 2 13.766 2 12.9622V10.9527C2 10.4169 2.41304 9.34515 4.06522 9.34515H11.5H18.9348C20.587 9.34515 21 9.98818 21 10.9527C21 11.9173 21 12.5603 21 12.9622C21 13.766 21 14.5697 19.7609 14.9716M3.23913 14.9716V20.1962C3.23913 20.4823 3.40435 21 4.06522 21C4.72609 21 14.2536 21 18.9348 21C19.2101 21 19.7609 20.8392 19.7609 20.1962C19.7609 19.5532 19.7609 14.9716 19.7609 14.9716M11.5 6.22459C12.7391 6.22459 15.2174 3.41135 11.5 1C7.78261 3.41135 10.2609 6.22459 11.5 6.22459Z'
      stroke='black'
      stroke-width='2'
      stroke-linejoin='round'
    />
  </svg>
);
