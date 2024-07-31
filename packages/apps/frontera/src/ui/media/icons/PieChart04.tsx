import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const PieChart04 = ({ className, ...props }: IconProps) => (
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
      d='M12 12L2.33178 9.44554C1.70974 11.7999 1.96813 14.3002 3.05847 16.4776C4.14881 18.6549 5.99616 20.3596 8.25395 21.2718L12 12ZM12 12L12.1047 2.00055C9.88558 1.97731 7.72174 2.69293 5.95402 4.0347C4.1863 5.37647 2.91514 7.26815 2.34075 9.41181L12 12ZM22 12C22 17.5228 17.5229 22 12 22C6.47716 22 2.00001 17.5228 2.00001 12C2.00001 6.47715 6.47716 2 12 2C17.5229 2 22 6.47715 22 12Z'
    />
  </svg>
);
