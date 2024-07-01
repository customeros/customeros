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
      d='M3.23913 15.9716C5.71739 16.7754 7.36957 15.9716 8.19565 13.9622C9.02174 15.9716 10.2609 16.3735 11.5 16.3735C13.1522 16.3735 13.9783 15.5697 14.8043 13.9622C15.2174 16.3735 18.9348 16.7754 19.7609 15.9716M3.23913 15.9716C2.54535 15.5216 2 14.766 2 13.9622V11.9527C2 11.4169 2.41304 10.3452 4.06522 10.3452H11.5H18.9348C20.587 10.3452 21 10.9882 21 11.9527C21 12.9173 21 13.5603 21 13.9622C21 14.766 21 15.5697 19.7609 15.9716M3.23913 15.9716V21.1962C3.23913 21.4823 3.40435 22 4.06522 22C4.72609 22 14.2536 22 18.9348 22C19.2101 22 19.7609 21.8392 19.7609 21.1962C19.7609 20.5532 19.7609 15.9716 19.7609 15.9716M11.5 7.22459C12.7391 7.22459 15.2174 4.41135 11.5 2C7.78261 4.41135 10.2609 7.22459 11.5 7.22459Z'
      stroke='#0C111D'
      stroke-width='2'
      stroke-linejoin='round'
    />
  </svg>
);
