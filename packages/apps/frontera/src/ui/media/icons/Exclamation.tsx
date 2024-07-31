import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Exclamation = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='3'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M12 6V11.9999M12 17.9999H12.01'
    />
  </svg>
);
