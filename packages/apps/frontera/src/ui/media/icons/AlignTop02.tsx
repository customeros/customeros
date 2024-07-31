import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const AlignTop02 = ({ className, ...props }: IconProps) => (
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
      d='M10 18V6C10 5.06812 10 4.60218 9.84776 4.23463C9.64477 3.74458 9.25542 3.35523 8.76537 3.15224C8.39782 3 7.93188 3 7 3C6.06812 3 5.60218 3 5.23463 3.15224C4.74458 3.35523 4.35523 3.74458 4.15224 4.23463C4 4.60218 4 5.06812 4 6V18C4 18.9319 4 19.3978 4.15224 19.7654C4.35523 20.2554 4.74458 20.6448 5.23463 20.8478C5.60218 21 6.06812 21 7 21C7.93188 21 8.39782 21 8.76537 20.8478C9.25542 20.6448 9.64477 20.2554 9.84776 19.7654C10 19.3978 10 18.9319 10 18Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M20 14V6C20 5.06812 20 4.60218 19.8478 4.23463C19.6448 3.74458 19.2554 3.35523 18.7654 3.15224C18.3978 3 17.9319 3 17 3C16.0681 3 15.6022 3 15.2346 3.15224C14.7446 3.35523 14.3552 3.74458 14.1522 4.23463C14 4.60218 14 5.06812 14 6V14C14 14.9319 14 15.3978 14.1522 15.7654C14.3552 16.2554 14.7446 16.6448 15.2346 16.8478C15.6022 17 16.0681 17 17 17C17.9319 17 18.3978 17 18.7654 16.8478C19.2554 16.6448 19.6448 16.2554 19.8478 15.7654C20 15.3978 20 14.9319 20 14Z'
    />
  </svg>
);
