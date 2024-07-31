import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ma = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#c1272d' d='M640 0H0v480h640z' />
    <path
      fill='none'
      stroke='#006233'
      strokeWidth='11.7'
      d='M320 179.4 284.4 289l93.2-67.6H262.4l93.2 67.6z'
    />
  </svg>
);
