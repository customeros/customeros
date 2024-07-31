import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CreditCardEdit = ({ className, ...props }: IconProps) => (
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
      d='M2 10H22V8.2C22 7.0799 22 6.51984 21.782 6.09202C21.5903 5.7157 21.2843 5.40974 20.908 5.21799C20.4802 5 19.9201 5 18.8 5H5.2C4.0799 5 3.51984 5 3.09202 5.21799C2.7157 5.40973 2.40973 5.71569 2.21799 6.09202C2 6.51984 2 7.0799 2 8.2V15.8C2 16.9201 2 17.4802 2.21799 17.908C2.40973 18.2843 2.71569 18.5903 3.09202 18.782C3.51984 19 4.0799 19 5.2 19H11M14.5 21L16.525 20.595C16.7015 20.5597 16.7898 20.542 16.8721 20.5097C16.9452 20.4811 17.0147 20.4439 17.079 20.399C17.1516 20.3484 17.2152 20.2848 17.3426 20.1574L21.5 16C22.0523 15.4477 22.0523 14.5523 21.5 14C20.9477 13.4477 20.0523 13.4477 19.5 14L15.3426 18.1574C15.2152 18.2848 15.1516 18.3484 15.101 18.421C15.0561 18.4853 15.0189 18.5548 14.9903 18.6279C14.958 18.7102 14.9403 18.7985 14.905 18.975L14.5 21Z'
    />
  </svg>
);
