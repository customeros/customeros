import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Twitter = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 32 32'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='currentColor'
      d='M3.89526 4L13.6761 17.2374L3.8335 28H6.04867L14.6658 18.5772L21.6283 28H29.1666L18.8355 14.018L27.9969 4H25.7817L17.8457 12.6782L11.4336 4H3.89526ZM7.15284 5.6516H10.616L25.9085 26.3481H22.4454L7.15284 5.6516Z'
    />
  </svg>
);
