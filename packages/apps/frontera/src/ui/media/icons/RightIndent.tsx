import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const RightIndent = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 16 16'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M14 2.66667H2M14 13.3333H2M8 6.16667H2M8 9.83334H2M13.1467 5.64L10.5689 7.57334C10.3759 7.71808 10.2794 7.79046 10.2449 7.87918C10.2147 7.95689 10.2147 8.04311 10.2449 8.12083C10.2794 8.20955 10.3759 8.28193 10.5689 8.42667L13.1467 10.36C13.4213 10.566 13.5586 10.669 13.6736 10.6666C13.7736 10.6645 13.8674 10.6176 13.9291 10.5388C14 10.4483 14 10.2767 14 9.93334V6.06667C14 5.72336 14 5.5517 13.9291 5.46117C13.8674 5.38239 13.7736 5.33549 13.6736 5.33341C13.5586 5.33102 13.4213 5.43402 13.1467 5.64Z'
    />
  </svg>
);
