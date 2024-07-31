import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Mail01 = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 12 10'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M1 2.5L5.08246 5.35772C5.41304 5.58913 5.57834 5.70484 5.75813 5.74965C5.91695 5.78924 6.08305 5.78924 6.24187 5.74965C6.42166 5.70484 6.58695 5.58913 6.91754 5.35772L11 2.5M3.4 9H8.6C9.44008 9 9.86012 9 10.181 8.83651C10.4632 8.6927 10.6927 8.46323 10.8365 8.18099C11 7.86012 11 7.44008 11 6.6V3.4C11 2.55992 11 2.13988 10.8365 1.81901C10.6927 1.53677 10.4632 1.3073 10.181 1.16349C9.86012 1 9.44008 1 8.6 1H3.4C2.55992 1 2.13988 1 1.81901 1.16349C1.53677 1.3073 1.3073 1.53677 1.16349 1.81901C1 2.13988 1 2.55992 1 3.4V6.6C1 7.44008 1 7.86012 1.16349 8.18099C1.3073 8.46323 1.53677 8.6927 1.81901 8.83651C2.13988 9 2.55992 9 3.4 9Z'
    />
  </svg>
);
