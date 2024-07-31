import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Underline02 = ({ className, ...props }: IconProps) => (
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
      d='M19 4V10C19 13.866 15.866 17 12 17C8.13401 17 5 13.866 5 10V4M8.5 4V10C8.5 13.2218 10.6766 15.9352 13.6395 16.7501M4 21H20M3 4L10.5 4M17 4L21 4'
    />
  </svg>
);
