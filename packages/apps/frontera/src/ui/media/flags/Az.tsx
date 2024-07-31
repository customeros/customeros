import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Az = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#3f9c35' d='M.1 0h640v480H.1z' />
    <path fill='#ed2939' d='M.1 0h640v320H.1z' />
    <path fill='#00b9e4' d='M.1 0h640v160H.1z' />
    <circle r='72' cx='304' cy='240' fill='#fff' />
    <circle r='60' cx='320' cy='240' fill='#ed2939' />
    <path
      fill='#fff'
      d='m384 200 7.7 21.5 20.6-9.8-9.8 20.7L424 240l-21.5 7.7 9.8 20.6-20.6-9.8L384 280l-7.7-21.5-20.6 9.8 9.8-20.6L344 240l21.5-7.7-9.8-20.6 20.6 9.8z'
    />
  </svg>
);
