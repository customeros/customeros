import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Xx = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='#fff'
      stroke='#adb5bd'
      strokeWidth='1.1'
      fillRule='evenodd'
      d='M.5.5h638.9v478.9H.5z'
    />
    <path
      fill='none'
      stroke='#adb5bd'
      strokeWidth='1.1'
      d='m.5.5 639 479m0-479-639 479'
    />
  </svg>
);
