import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const MessageTextSquare01 = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 10 11'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      d='M2.5 3.25H5M2.5 5H6.5M3.84187 8H7.1C7.94008 8 8.36012 8 8.68099 7.83651C8.96323 7.6927 9.1927 7.46323 9.33651 7.18099C9.5 6.86012 9.5 6.44008 9.5 5.6V2.9C9.5 2.05992 9.5 1.63988 9.33651 1.31901C9.1927 1.03677 8.96323 0.8073 8.68099 0.66349C8.36012 0.5 7.94008 0.5 7.1 0.5H2.9C2.05992 0.5 1.63988 0.5 1.31901 0.66349C1.03677 0.8073 0.8073 1.03677 0.66349 1.31901C0.5 1.63988 0.5 2.05992 0.5 2.9V9.16775C0.5 9.43418 0.5 9.56739 0.554614 9.63581C0.602111 9.69531 0.674135 9.72993 0.75027 9.72984C0.837814 9.72975 0.941835 9.64653 1.14988 9.4801L2.34261 8.52592C2.58626 8.33099 2.70808 8.23353 2.84374 8.16423C2.9641 8.10274 3.09221 8.0578 3.22461 8.03063C3.37383 8 3.52985 8 3.84187 8Z'
      stroke='#667085'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
