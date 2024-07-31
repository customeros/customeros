import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const HorizontalBarChart03 = ({ className, ...props }: IconProps) => (
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
      d='M13 9.5V6.1C13 5.53995 13 5.25992 12.891 5.04601C12.7951 4.85785 12.6422 4.70487 12.454 4.60899C12.2401 4.5 11.9601 4.5 11.4 4.5H3M17 14.5V11.1C17 10.5399 17 10.2599 16.891 10.046C16.7951 9.85785 16.6422 9.70487 16.454 9.60899C16.2401 9.5 15.9601 9.5 15.4 9.5H3M3 2L3 22M3 19.5H19.4C19.9601 19.5 20.2401 19.5 20.454 19.391C20.6422 19.2951 20.7951 19.1422 20.891 18.954C21 18.7401 21 18.4601 21 17.9V16.1C21 15.5399 21 15.2599 20.891 15.046C20.7951 14.8578 20.6422 14.7049 20.454 14.609C20.2401 14.5 19.9601 14.5 19.4 14.5L3 14.5L3 19.5Z'
    />
  </svg>
);
