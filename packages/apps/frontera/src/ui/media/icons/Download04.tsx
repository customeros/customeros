import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Download04 = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      d='M8 12L12 16M12 16L16 12M12 16V6.8C12 5.40929 12 4.71394 11.4495 3.9354C11.0837 3.41812 10.0306 2.77968 9.40278 2.69462C8.45789 2.5666 8.09907 2.75378 7.38143 3.12814C4.18333 4.79643 2 8.14324 2 12C2 17.5228 6.47715 22 12 22C17.5228 22 22 17.5228 22 12C22 8.29859 19.989 5.06687 17 3.33782'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
