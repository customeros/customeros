import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Fingerprint02 = ({ className, ...props }: IconProps) => (
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
      d='M12 10V14M7.44712 3.42103C8.73941 2.52503 10.3084 2 12 2C16.4183 2 20 5.58172 20 10V11.2367M4.41632 7.44607C4.14633 8.24809 4 9.10696 4 10V14C4 17.6349 6.42416 20.7035 9.74396 21.6775M19.6588 16.3187C18.9294 18.7314 17.0911 20.6626 14.7367 21.5196M14.325 6.14635C13.6464 5.7361 12.8508 5.5 12 5.5C9.51472 5.5 7.5 7.51472 7.5 10V12.95M16.5 11.04V14C16.5 16.4853 14.4853 18.5 12 18.5C11.1514 18.5 10.3576 18.2651 9.68014 17.8567'
    />
  </svg>
);
