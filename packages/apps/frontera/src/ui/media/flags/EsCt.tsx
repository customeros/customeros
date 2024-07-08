import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const EsCt = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fcdd09' d='M0 0h640v480H0z' />
    <path
      stroke='#da121a'
      strokeWidth='60'
      d='M0 90h810m0 120H0m0 120h810m0 120H0'
      transform='scale(.79012 .88889)'
    />
  </svg>
);
