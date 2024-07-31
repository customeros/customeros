import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Tg = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='tg-a'>
        <path fillOpacity='.7' d='M0 0h682.7v512H0z' />
      </clipPath>
    </defs>
    <g fillRule='evenodd' clipPath='url(#tg-a)' transform='scale(.9375)'>
      <path fill='#ffe300' d='M0 0h767.6v512H0z' />
      <path fill='#118600' d='M0 208.1h767.6V311H0zM0 .2h767.6v102.9H0z' />
      <path fill='#d80000' d='M0 .3h306.5v310.6H0z' />
      <path
        fill='#fff'
        d='M134.4 128.4c0-.8 18.9-53 18.9-53l17 52.2s57.4 1.7 57.4.8-45.3 34.3-45.3 34.3 21.4 60 20.5 58.2c-.8-1.7-49.6-36-49.6-36s-49.7 34.3-48.8 34.3c.8 0 18.8-56.5 18.8-56.5l-44.5-33.4z'
      />
      <path fill='#118600' d='M0 409.2h767.6V512H0z' />
    </g>
  </svg>
);
