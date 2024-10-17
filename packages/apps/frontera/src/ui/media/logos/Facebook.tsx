import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Facebook = ({ className, ...props }: IconProps) => (
  <svg
    width='32'
    height='32'
    fill='none'
    viewBox='0 0 32 32'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <circle r='14' cx='16' cy='16' fill='url(#paint0_linear_1334_718)' />
    <path
      fill='white'
      d='M21.2137 20.2816L21.8356 16.3301H17.9452V13.767C17.9452 12.6857 18.4877 11.6311 20.2302 11.6311H22V8.26699C22 8.26699 20.3945 8 18.8603 8C15.6548 8 13.5617 9.89294 13.5617 13.3184V16.3301H10V20.2816H13.5617V29.8345C14.2767 29.944 15.0082 30 15.7534 30C16.4986 30 17.2302 29.944 17.9452 29.8345V20.2816H21.2137Z'
    />
    <defs>
      <linearGradient
        y1='2'
        x1='16'
        x2='16'
        y2='29.917'
        id='paint0_linear_1334_718'
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#18ACFE' />
        <stop offset='1' stopColor='#0163E0' />
      </linearGradient>
    </defs>
  </svg>
);
