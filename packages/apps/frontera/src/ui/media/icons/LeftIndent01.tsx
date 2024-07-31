import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LeftIndent01 = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M21 9.25H12M21 4H3M21 14.75H12M21 20H3M4.28 8.56L8.14667 11.46C8.43616 11.6771 8.5809 11.7857 8.63266 11.9188C8.678 12.0353 8.678 12.1647 8.63266 12.2812C8.5809 12.4143 8.43616 12.5229 8.14667 12.74L4.28 15.64C3.86802 15.949 3.66203 16.1035 3.48961 16.0999C3.33956 16.0968 3.19885 16.0264 3.10632 15.9082C3 15.7725 3 15.515 3 15V9.2C3 8.68503 3 8.42754 3.10632 8.29175C3.19885 8.17358 3.33956 8.10323 3.48961 8.10011C3.66203 8.09652 3.86802 8.25102 4.28 8.56Z'
    />
  </svg>
);
