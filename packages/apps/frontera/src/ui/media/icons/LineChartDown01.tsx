import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LineChartDown01 = ({ className, ...props }: IconProps) => (
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
      d='M21 21H4.6C4.03995 21 3.75992 21 3.54601 20.891C3.35785 20.7951 3.20487 20.6422 3.10899 20.454C3 20.2401 3 19.9601 3 19.4V3M20 15L16.0811 10.8173C15.9326 10.6588 15.8584 10.5796 15.7688 10.5386C15.6897 10.5024 15.6026 10.4875 15.516 10.4953C15.4179 10.5042 15.3215 10.5542 15.1287 10.6543L11.8713 12.3457C11.6785 12.4458 11.5821 12.4958 11.484 12.5047C11.3974 12.5125 11.3103 12.4976 11.2312 12.4614C11.1416 12.4204 11.0674 12.3412 10.9189 12.1827L7 8'
    />
  </svg>
);
