import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sb = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='sb-a'>
        <path fillOpacity='.7' d='M0 0h682.7v512H0z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      strokeWidth='1pt'
      clipPath='url(#sb-a)'
      transform='scale(.9375)'
    >
      <path fill='#0000d6' d='M0 507.2 987.4 0H0z' />
      <path fill='#006000' d='M1024 0 27.2 512H1024z' />
      <path fill='#fc0' d='M1024 0h-54.9L0 485.4V512h54.9L1024 27.6z' />
      <path
        fill='#fff'
        d='m71.4 9.1 11.8 34.5h38.5L90.6 64.7l11.9 34.4L71.4 78 40.3 99.2l11.9-34.4-31.1-21.3h38.4zm191.1 0 11.9 34.5h38.5l-31.2 21.2 12 34.4L262.4 78l-31 21.3 11.9-34.4-31.2-21.3h38.5zm0 144.5 11.9 34.5h38.5l-31.2 21.2 12 34.4-31.2-21.3-31 21.3 11.9-34.4-31.2-21.3h38.5zm-95-71.4 11.9 34.4h38.4l-31 21.3 11.8 34.4-31-21.3-31.2 21.3 12-34.4-31.2-21.3h38.5zm-96.1 71.4 11.8 34.5h38.5l-31.1 21.2 11.9 34.4-31.1-21.3-31.1 21.3 12-34.4L21 188h38.4z'
      />
    </g>
  </svg>
);