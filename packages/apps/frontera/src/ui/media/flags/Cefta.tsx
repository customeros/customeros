import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cefta = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#039' d='M0 0h640v480H0z' />
    <circle
      cx='320'
      r='30.4'
      cy='249.8'
      fill='none'
      stroke='#fc0'
      strokeWidth='27.5'
    />
    <circle
      cx='320'
      r='88.3'
      cy='249.8'
      fill='none'
      stroke='#fc0'
      strokeWidth='27.5'
    />
    <path fill='#039' d='m404.7 165.1 84.7 84.7-84.7 84.7-84.7-84.7z' />
    <path
      fill='#fc0'
      d='M175.7 236.1h59.2v27.5h-59.2zm259.1 0h88.3v27.5h-88.3zM363 187.4l38.8-38.8 19.4 19.5-38.7 38.7zM306.3 48.6h27.5v107.1h-27.5z'
    />
    <circle r='13.7' cx='225.1' cy='159.6' fill='#fc0' />
    <circle r='13.7' cx='144.3' cy='249.8' fill='#fc0' />
    <circle cx='320' r='13.7' cy='381.4' fill='#fc0' />
    <circle cx='320' r='13.7' cy='425.5' fill='#fc0' />
    <circle r='13.7' cx='408.3' cy='249.8' fill='#fc0' />
    <path
      fill='#fc0'
      d='m208.3 341.5 19.5-19.4 19.4 19.4-19.4 19.5zm204.7 21 19.5-19.5 19.5 19.5-19.5 19.4z'
    />
  </svg>
);
