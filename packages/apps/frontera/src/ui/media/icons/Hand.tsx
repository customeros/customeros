import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hand = ({ className, ...props }: IconProps) => (
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
      d='M6.9 11.4444V14.2222M6.9 11.4444V4.77778C6.9 3.8573 7.66112 3.11111 8.6 3.11111C9.53888 3.11111 10.3 3.8573 10.3 4.77778M6.9 11.4444C6.9 10.524 6.13888 9.77778 5.2 9.77778C4.26112 9.77778 3.5 10.524 3.5 11.4444V13.6667C3.5 18.269 7.30558 22 12 22C16.6944 22 20.5 18.269 20.5 13.6667V8.11111C20.5 7.19064 19.7389 6.44444 18.8 6.44444C17.8611 6.44444 17.1 7.19064 17.1 8.11111M10.3 4.77778V10.8889M10.3 4.77778V3.66667C10.3 2.74619 11.0611 2 12 2C12.9389 2 13.7 2.74619 13.7 3.66667V4.77778M13.7 4.77778V10.8889M13.7 4.77778C13.7 3.8573 14.4611 3.11111 15.4 3.11111C16.3389 3.11111 17.1 3.8573 17.1 4.77778V8.11111M17.1 8.11111V10.8889'
    />
  </svg>
);
