import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hourglass02 = ({ className, ...props }: IconProps) => (
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
      d='M18.1626 2H5.83744C5.37494 2 5 2.37494 5 2.83744C5 5.50268 6.05876 8.05876 7.94337 9.94337L9.16256 11.1626C9.28363 11.2836 9.34417 11.3442 9.3875 11.4023C9.65188 11.7569 9.65188 12.2431 9.3875 12.5977C9.34417 12.6558 9.28363 12.7164 9.16256 12.8374L7.94337 14.0566C6.05876 15.9412 5 18.4973 5 21.1626C5 21.6251 5.37494 22 5.83744 22H18.1626C18.6251 22 19 21.6251 19 21.1626C19 18.4973 17.9412 15.9412 16.0566 14.0566L14.8374 12.8374C14.7164 12.7164 14.6558 12.6558 14.6125 12.5977C14.3481 12.2431 14.3481 11.7569 14.6125 11.4023C14.6558 11.3442 14.7164 11.2836 14.8374 11.1626L16.0566 9.94337C17.9412 8.05876 19 5.50268 19 2.83744C19 2.37494 18.6251 2 18.1626 2Z'
    />
  </svg>
);
