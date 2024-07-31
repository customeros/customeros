import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Airpods = ({ className, ...props }: IconProps) => (
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
      d='M2 7.625C2 9.90317 3.84683 11.75 6.125 11.75C6.43089 11.75 6.58383 11.75 6.66308 11.7773C6.82888 11.8345 6.91545 11.9211 6.97266 12.0869C7 12.1662 7 12.2903 7 12.5386V18.875C7 19.7725 7.72754 20.5 8.625 20.5C9.52246 20.5 10.25 19.7725 10.25 18.875V7.625C10.25 5.34683 8.40317 3.5 6.125 3.5C3.84683 3.5 2 5.34683 2 7.625Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M22 7.625C22 9.90317 20.1532 11.75 17.875 11.75C17.5691 11.75 17.4162 11.75 17.3369 11.7773C17.1711 11.8345 17.0845 11.9211 17.0273 12.0869C17 12.1662 17 12.2903 17 12.5386V18.875C17 19.7725 16.2725 20.5 15.375 20.5C14.4775 20.5 13.75 19.7725 13.75 18.875V7.625C13.75 5.34683 15.5968 3.5 17.875 3.5C20.1532 3.5 22 5.34683 22 7.625Z'
    />
  </svg>
);
