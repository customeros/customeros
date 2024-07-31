import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ChartBreakoutCircle = ({ className, ...props }: IconProps) => (
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
      d='M15.5 3.5V2M19.4393 4.56066L20.5 3.5M20.5103 8.5H22.0103M21.9506 13C21.4489 18.0533 17.1853 22 12 22C6.47715 22 2 17.5228 2 12C2 6.81465 5.94668 2.5511 11 2.04938M12 8H16V12M15.6197 8C13.2653 11.3276 9.38636 13.5 5 13.5C3.9971 13.5 3.02072 13.3864 2.08302 13.1715'
    />
  </svg>
);
