import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Atom02 = ({ className, ...props }: IconProps) => (
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
      d='M17.1153 15.3582C16.8446 15.6642 16.5606 15.9665 16.2635 16.2635C11.9678 20.5593 6.58585 22.1422 4.2427 19.7991C2.6363 18.1926 2.8752 15.158 4.56847 12.0242M6.88967 8.72526C7.17138 8.40495 7.46772 8.08875 7.77824 7.77824C12.074 3.48247 17.4559 1.89956 19.7991 4.2427C21.4066 5.85021 21.1662 8.88795 19.4698 12.024M16.2635 7.77824C20.5593 12.074 22.1422 17.4559 19.7991 19.7991C17.4559 22.1422 12.074 20.5593 7.77824 16.2635C3.48247 11.9678 1.89956 6.58585 4.2427 4.2427C6.58585 1.89956 11.9678 3.48247 16.2635 7.77824ZM13.0001 12C13.0001 12.5523 12.5523 13 12.0001 13C11.4478 13 11.0001 12.5523 11.0001 12C11.0001 11.4477 11.4478 11 12.0001 11C12.5523 11 13.0001 11.4477 13.0001 12Z'
    />
  </svg>
);
