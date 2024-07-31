import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Backpack = ({ className, ...props }: IconProps) => (
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
      d='M20 13V17.8C20 18.9201 20 19.4802 19.782 19.908C19.5903 20.2843 19.2843 20.5903 18.908 20.782C18.4802 21 17.9201 21 16.8 21H7.2C6.0799 21 5.51984 21 5.09202 20.782C4.71569 20.5903 4.40973 20.2843 4.21799 19.908C4 19.4802 4 18.9201 4 17.8V13M9 10H15M9.28571 14H14.7143C16.8467 14 17.913 14 18.7355 13.6039C19.552 13.2107 20.2107 12.552 20.6039 11.7355C21 10.913 21 9.84674 21 7.71429C21 6.11494 21 5.31527 20.7029 4.69835C20.408 4.08603 19.914 3.59197 19.3017 3.29709C18.6847 3 17.8851 3 16.2857 3H7.71429C6.11494 3 5.31527 3 4.69835 3.29709C4.08603 3.59197 3.59197 4.08603 3.29709 4.69835C3 5.31527 3 6.11494 3 7.71429C3 9.84674 3 10.913 3.39612 11.7355C3.7893 12.552 4.44803 13.2107 5.26447 13.6039C6.08703 14 7.15326 14 9.28571 14Z'
    />
  </svg>
);
