import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Underline = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 16 16'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M12 2.66667V7.33334C12 9.54248 10.2091 11.3333 8 11.3333C5.79086 11.3333 4 9.54248 4 7.33334V2.66667M2.66667 14H13.3333'
    />
  </svg>
);
