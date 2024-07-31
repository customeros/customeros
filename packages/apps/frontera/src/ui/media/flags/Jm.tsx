import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Jm = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#000001' d='m0 0 320 240L0 480zm640 0L320 240l320 240z' />
      <path fill='#090' d='m0 0 320 240L640 0zm0 480 320-240 320 240z' />
      <path fill='#fc0' d='M640 0h-59.6L0 435.3V480h59.6L640 44.7z' />
      <path fill='#fc0' d='M0 0v44.7L580.4 480H640v-44.7L59.6 0z' />
    </g>
  </svg>
);
