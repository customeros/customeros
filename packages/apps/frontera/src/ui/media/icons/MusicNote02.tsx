import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const MusicNote02 = ({ className, ...props }: IconProps) => (
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
      d='M12 18V5.58888C12 4.73166 12 4.30306 12.1805 4.04492C12.3382 3.81952 12.5817 3.669 12.8538 3.6288C13.1655 3.58276 13.5488 3.77444 14.3155 4.1578L18 6.00003M12 18C12 19.6569 10.6569 21 9 21C7.34315 21 6 19.6569 6 18C6 16.3432 7.34315 15 9 15C10.6569 15 12 16.3432 12 18Z'
    />
  </svg>
);
