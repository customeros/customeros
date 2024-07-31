import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ArrowDropdown = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='currentColor'
      d='M14.7929 10H9.20711C8.76165 10 8.53857 10.5386 8.85355 10.8536L11.6464 13.6464C11.8417 13.8417 12.1583 13.8417 12.3536 13.6464L15.1464 10.8536C15.4614 10.5386 15.2383 10 14.7929 10Z'
    />
  </svg>
);
