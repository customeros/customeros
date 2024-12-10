import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const GitTimeline = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      stroke-width='2'
      stroke='currentColor'
      stroke-linecap='round'
      stroke-linejoin='round'
      d='M12 16C14.2091 16 16 14.2091 16 12C16 9.79086 14.2091 8 12 8M12 16C9.79086 16 8 14.2091 8 12C8 9.79086 9.79086 8 12 8M12 16V22M12 8L12 2'
    />
  </svg>
);
