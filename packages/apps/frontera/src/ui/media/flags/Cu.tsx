import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cu = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='cu-a'>
        <path fillOpacity='.7' d='M-32 0h682.7v512H-32z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      clipPath='url(#cu-a)'
      transform='translate(30)scale(.94)'
    >
      <path fill='#002a8f' d='M-32 0h768v512H-32z' />
      <path fill='#fff' d='M-32 102.4h768v102.4H-32zm0 204.8h768v102.4H-32z' />
      <path fill='#cb1515' d='m-32 0 440.7 255.7L-32 511z' />
      <path
        fill='#fff'
        d='M161.8 325.5 114.3 290l-47.2 35.8 17.6-58.1-47.2-36 58.3-.4 18.1-58 18.5 57.8 58.3.1-46.9 36.3z'
      />
    </g>
  </svg>
);
