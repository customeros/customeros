import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const AlignBottom02 = ({ className, ...props }: IconProps) => (
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
      d='M20 18V10C20 9.06812 20 8.60218 19.8478 8.23463C19.6448 7.74458 19.2554 7.35523 18.7654 7.15224C18.3978 7 17.9319 7 17 7C16.0681 7 15.6022 7 15.2346 7.15224C14.7446 7.35523 14.3552 7.74458 14.1522 8.23463C14 8.60218 14 9.06812 14 10V18C14 18.9319 14 19.3978 14.1522 19.7654C14.3552 20.2554 14.7446 20.6448 15.2346 20.8478C15.6022 21 16.0681 21 17 21C17.9319 21 18.3978 21 18.7654 20.8478C19.2554 20.6448 19.6448 20.2554 19.8478 19.7654C20 19.3978 20 18.9319 20 18Z'
    />
  </svg>
);
