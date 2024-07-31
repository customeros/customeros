import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Compass02 = ({ className, ...props }: IconProps) => (
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
      d='M22 12H20M22 12C22 17.5228 17.5228 22 12 22M22 12C22 6.47715 17.5228 2 12 2M19.0711 19.0711L17.6569 17.6569M4 12H2M2 12C2 17.5228 6.47715 22 12 22M2 12C2 6.47715 6.47715 2 12 2M6.34315 6.34315L4.92893 4.92893M12 4V2M17.6569 6.34315L19.0711 4.92893M12 22V20M4.92893 19.0711L6.34315 17.6569M12 8L16 12L12 16L8 12L12 8Z'
    />
  </svg>
);
