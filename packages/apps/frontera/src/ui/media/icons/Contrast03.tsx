import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Contrast03 = ({ className, ...props }: IconProps) => (
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
      d='M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M16 8.5C16 12.6421 12.6421 16 8.5 16C7.88534 16 7.28795 15.9261 6.71623 15.7866C7.89585 17.4297 9.82294 18.5 12 18.5C15.5899 18.5 18.5 15.5899 18.5 12C18.5 9.82294 17.4297 7.89585 15.7866 6.71623C15.9261 7.28795 16 7.88534 16 8.5Z'
    />
  </svg>
);
