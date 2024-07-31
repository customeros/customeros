import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Glasses01 = ({ className, ...props }: IconProps) => (
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
      d='M10 11.5347C11.2335 10.8218 12.7663 10.8218 13.9999 11.5347M8.82843 9.17157C10.3905 10.7337 10.3905 13.2663 8.82843 14.8284C7.26634 16.3905 4.73367 16.3905 3.17157 14.8284C1.60948 13.2663 1.60948 10.7337 3.17157 9.17157C4.73366 7.60948 7.26633 7.60948 8.82843 9.17157ZM20.8284 9.17157C22.3905 10.7337 22.3905 13.2663 20.8284 14.8284C19.2663 16.3905 16.7337 16.3905 15.1716 14.8284C13.6095 13.2663 13.6095 10.7337 15.1716 9.17157C16.7337 7.60948 19.2663 7.60948 20.8284 9.17157Z'
    />
  </svg>
);
