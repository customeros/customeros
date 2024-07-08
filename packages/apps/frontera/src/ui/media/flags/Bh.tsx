import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bh = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0' />
    <path
      fill='#ce1126'
      d='M640 0H96l110.7 48L96 96l110.7 48L96 192l110.7 48L96 288l110.7 48L96 384l110.7 48L96 480h544'
    />
  </svg>
);
