import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const WatchCircle = ({ className, ...props }: IconProps) => (
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
      d='M7 17L7.48551 19.4276C7.66878 20.3439 7.76041 20.8021 7.99964 21.1448C8.2106 21.447 8.50097 21.685 8.83869 21.8326C9.22166 22 9.6889 22 10.6234 22H13.3766C14.3111 22 14.7783 22 15.1613 21.8326C15.499 21.685 15.7894 21.447 16.0004 21.1448C16.2396 20.8021 16.3312 20.3439 16.5145 19.4276L17 17M7 7L7.48551 4.57243C7.66878 3.6561 7.76041 3.19793 7.99964 2.85522C8.2106 2.55301 8.50097 2.31497 8.83869 2.16737C9.22166 2 9.6889 2 10.6234 2H13.3766C14.3111 2 14.7783 2 15.1613 2.16737C15.499 2.31497 15.7894 2.55301 16.0004 2.85522C16.2396 3.19793 16.3312 3.6561 16.5145 4.57243L17 7M12 9V12L13.5 13.5M19 12C19 15.866 15.866 19 12 19C8.13401 19 5 15.866 5 12C5 8.13401 8.13401 5 12 5C15.866 5 19 8.13401 19 12Z'
    />
  </svg>
);
