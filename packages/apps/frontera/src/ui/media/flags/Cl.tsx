import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cl = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='cl-a'>
        <path fillOpacity='.7' d='M0 0h682.7v512H0z' />
      </clipPath>
    </defs>
    <g fillRule='evenodd' clipPath='url(#cl-a)' transform='scale(.9375)'>
      <path fill='#fff' d='M256 0h512v256H256z' />
      <path fill='#0039a6' d='M0 0h256v256H0z' />
      <path
        fill='#fff'
        d='M167.8 191.7 128.2 162l-39.5 30 14.7-48.8L64 113.1l48.7-.5L127.8 64l15.5 48.5 48.7.1-39.2 30.4z'
      />
      <path fill='#d52b1e' d='M0 256h768v256H0z' />
    </g>
  </svg>
);
