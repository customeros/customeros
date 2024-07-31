import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bs = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='bs-a'>
        <path fillOpacity='.7' d='M-12 0h640v480H-12z' />
      </clipPath>
    </defs>
    <g fillRule='evenodd' clipPath='url(#bs-a)' transform='translate(12)'>
      <path fill='#fff' d='M968.5 480h-979V1.8h979z' />
      <path fill='#ffe900' d='M968.5 344.5h-979V143.3h979z' />
      <path fill='#08ced6' d='M968.5 480h-979V320.6h979zm0-318.7h-979V2h979z' />
      <path
        fill='#000001'
        d='M-11 0c2.3 0 391.8 236.8 391.8 236.8L-12 479.2z'
      />
    </g>
  </svg>
);
