import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const GitCommit = ({ className, ...props }: IconProps) => (
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
      d='M12 8.0002C14.2091 8.0002 16 9.79106 16 12.0002C16 14.2093 14.2091 16.0002 12 16.0002M12 8.0002C9.79086 8.0002 8 9.79106 8 12.0002C8 14.2093 9.79086 16.0002 12 16.0002M12 8.0002L12 2.00019M12 16.0002L12 22'
    />
  </svg>
);
