import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Radar = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      strokeWidth={props?.strokeWidth || 2}
      d='M22 12c0 5.523-4.477 10-10 10a9.954 9.954 0 0 1-5.251-1.488M18 12A6 6 0 1 1 6.749 9.095m6.574 1.405a2 2 0 1 0-2.646 3 2 2 0 0 0 2.646-3Zm0 0 2.119-3.415m0 0A5.972 5.972 0 0 0 12 6a5.972 5.972 0 0 0-3.001.803m6.443.282 2.127-3.392A9.954 9.954 0 0 0 12 2C6.477 2 2 6.477 2 12c0 2.394.841 4.592 2.244 6.313'
    />
  </svg>
);
