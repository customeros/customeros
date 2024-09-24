import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LinkedinOutline = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <rect
      x='2.5'
      y='2.5'
      rx='9.5'
      width='19'
      height='19'
      strokeWidth='2'
      stroke='currentColor'
    />
    <path
      fill='currentColor'
      d='M9.46392 7.26911C9.46392 7.97002 8.85639 8.53822 8.10696 8.53822C7.35753 8.53822 6.75 7.97002 6.75 7.26911C6.75 6.5682 7.35753 6 8.10696 6C8.85639 6 9.46392 6.5682 9.46392 7.26911Z'
    />
    <path
      fill='currentColor'
      d='M6.93557 9.47107H9.25515V16.5H6.93557V9.47107Z'
    />
    <path
      fill='currentColor'
      d='M12.9897 9.47107H10.6701V16.5H12.9897C12.9897 16.5 12.9897 14.2872 12.9897 12.9036C12.9897 12.0732 13.2732 11.2392 14.4046 11.2392C15.6833 11.2392 15.6756 12.3259 15.6696 13.1678C15.6618 14.2683 15.6804 15.3914 15.6804 16.5H18V12.7903C17.9804 10.4215 17.3631 9.33006 15.3325 9.33006C14.1265 9.33006 13.379 9.87754 12.9897 10.3729V9.47107Z'
    />
  </svg>
);
