import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const So = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='so-a'>
        <path fillOpacity='.7' d='M-85.3 0h682.6v512H-85.3z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      clipPath='url(#so-a)'
      transform='translate(80)scale(.9375)'
    >
      <path fill='#40a6ff' d='M-128 0h768v512h-768z' />
      <path
        fill='#fff'
        d='M336.5 381.2 254 327.7l-82.1 54 30.5-87.7-82-54.2L222 239l31.4-87.5 32.1 87.3 101.4.1-81.5 54.7z'
      />
    </g>
  </svg>
);
