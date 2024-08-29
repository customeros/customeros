import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const BuildingAdd = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      stroke-width='2'
      stroke='currentColor'
      stroke-linecap='round'
      stroke-linejoin='round'
      d='M7.5 11H4.6C4.03995 11 3.75992 11 3.54601 11.109C3.35785 11.2049 3.20487 11.3578 3.10899 11.546C3 11.7599 3 12.0399 3 12.6V21M16.5 21H2M7.5 21V6.2C7.5 5.0799 7.5 4.51984 7.71799 4.09202C7.90973 3.71569 8.21569 3.40973 8.59202 3.21799C9.01984 3 9.57989 3 10.7 3H13.3C14.4201 3 14.9802 3 15.408 3.21799C15.7843 3.40973 16.0903 3.71569 16.282 4.09202C16.5 4.51984 16.5 5.0799 16.5 6.2V11M11 7H13M11 11H13M19.5 19V16M19.5 16V13M19.5 16H16.5M19.5 16H22.5'
    />
  </svg>
);
