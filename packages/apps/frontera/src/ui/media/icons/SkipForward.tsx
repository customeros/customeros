import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const SkipForward = ({ className, ...props }: IconProps) => (
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
      d='M19 5V19M7.59951 17.9204L13.4383 13.2494C13.9721 12.8223 14.239 12.6088 14.3357 12.3508C14.4204 12.1246 14.4204 11.8754 14.3357 11.6492C14.239 11.3912 13.9721 11.1777 13.4383 10.7506L7.59951 6.07961C6.76734 5.41387 6.35125 5.081 6.00108 5.08063C5.69654 5.0803 5.40845 5.21876 5.21846 5.45677C5 5.73045 5 6.2633 5 7.329V16.671C5 17.7367 5 18.2695 5.21846 18.5432C5.40845 18.7812 5.69654 18.9197 6.00108 18.9194C6.35125 18.919 6.76734 18.5861 7.59951 17.9204Z'
    />
  </svg>
);
