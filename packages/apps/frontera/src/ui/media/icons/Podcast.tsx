import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Podcast = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M17.1189 18C19.4623 16.4151 21 13.7779 21 10.785C21 5.9333 16.9704 2 12 2C7.02958 2 3 5.9333 3 10.785C3 13.7779 4.53771 16.4151 6.88113 18M8.35967 14C7.51875 13.15 7 12.0086 7 10.7505C7 8.12711 9.23881 6 12 6C14.7612 6 17 8.12711 17 10.7505C17 12.0095 16.4813 13.15 15.6403 14M12 22C10.8954 22 10 21.1046 10 20V18C10 16.8954 10.8954 16 12 16C13.1046 16 14 16.8954 14 18V20C14 21.1046 13.1046 22 12 22ZM13 11C13 11.5523 12.5523 12 12 12C11.4477 12 11 11.5523 11 11C11 10.4477 11.4477 10 12 10C12.5523 10 13 10.4477 13 11Z'
    />
  </svg>
);
