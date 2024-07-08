import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gh = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#006b3f' d='M0 0h640v480H0z' />
    <path fill='#fcd116' d='M0 0h640v320H0z' />
    <path fill='#ce1126' d='M0 0h640v160H0z' />
    <path fill='#000001' d='m320 160 52 160-136.1-98.9H404L268 320z' />
  </svg>
);
