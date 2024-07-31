import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LeftIndent = ({ className, ...props }: IconProps) => (
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
      d='M14 6.16667H8M14 2.66667H2M14 9.83334H8M14 13.3333H2M2.85333 5.70667L5.43111 7.64C5.62411 7.78475 5.7206 7.85712 5.75511 7.94585C5.78533 8.02356 5.78533 8.10978 5.75511 8.1875C5.7206 8.27622 5.62411 8.34859 5.43111 8.49334L2.85333 10.4267C2.57868 10.6327 2.44135 10.7357 2.3264 10.7333C2.22637 10.7312 2.13256 10.6843 2.07088 10.6055C2 10.515 2 10.3433 2 10V6.13334C2 5.79002 2 5.61836 2.07088 5.52784C2.13256 5.44906 2.22637 5.40216 2.3264 5.40008C2.44135 5.39769 2.57868 5.50068 2.85333 5.70667Z'
    />
  </svg>
);
