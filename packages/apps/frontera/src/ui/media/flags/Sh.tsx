import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sh = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#012169' d='M0 0h640v480H0z' />
    <path
      fill='#FFF'
      d='m75 0 244 181L562 0h78v62L400 241l240 178v61h-80L320 301 81 480H0v-60l239-178L0 64V0z'
    />
    <path
      fill='#C8102E'
      d='m424 281 216 159v40L369 281zm-184 20 6 35L54 480H0zM640 0v3L391 191l2-44L590 0zM0 0l239 176h-60L0 42z'
    />
    <path fill='#FFF' d='M241 0v480h160V0zM0 160v160h640V160z' />
    <path fill='#C8102E' d='M0 193v96h640v-96zM273 0v480h96V0z' />
  </svg>
);
