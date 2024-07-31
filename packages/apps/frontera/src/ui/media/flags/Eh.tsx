import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Eh = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='eh-a'>
        <path fillOpacity='.7' d='M-158.7 0H524v512h-682.7z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      clipPath='url(#eh-a)'
      transform='translate(148.8)scale(.94)'
    >
      <path fill='#000001' d='M-158.3 0h680.9v255.3h-680.9z' />
      <path fill='#007a3d' d='M-158.3 255.3h680.9v255.3h-680.9z' />
      <path fill='#fff' d='M-158.3 148.9h680.9v212.8h-680.9z' />
      <path fill='#c4111b' d='m-158.3 0 340.4 255.3-340.4 255.3Z' />
      <circle r='68.1' cx='352.3' cy='255.3' fill='#c4111b' />
      <circle r='68.1' cx='377.9' cy='255.3' fill='#fff' />
      <path
        fill='#c4111b'
        d='m334 296.5 29.1-20.7 28.8 21-10.8-34 29-20.9-35.7-.2-11-34-11.2 33.9-35.7-.2 28.7 21.2-11.1 34z'
      />
    </g>
  </svg>
);
