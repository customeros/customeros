import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const PaintPour = ({ className, ...props }: IconProps) => (
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
      d='M15.9997 11L1.9997 11M9.9997 4L7.9997 2M13.9997 22L1.9997 22M21.9997 16C21.9997 17.1046 21.1043 18 19.9997 18C18.8951 18 17.9997 17.1046 17.9997 16C17.9997 14.8954 19.9997 13 19.9997 13C19.9997 13 21.9997 14.8954 21.9997 16ZM8.9997 3L15.8683 9.86863C16.2643 10.2646 16.4624 10.4627 16.5365 10.691C16.6018 10.8918 16.6018 11.1082 16.5365 11.309C16.4624 11.5373 16.2643 11.7354 15.8683 12.1314L11.2624 16.7373C10.4704 17.5293 10.0744 17.9253 9.61773 18.0737C9.21605 18.2042 8.78335 18.2042 8.38166 18.0737C7.92501 17.9253 7.52899 17.5293 6.73696 16.7373L3.26244 13.2627C2.4704 12.4707 2.07439 12.0747 1.92601 11.618C1.7955 11.2163 1.7955 10.7837 1.92601 10.382C2.07439 9.92531 2.47041 9.52929 3.26244 8.73726L8.9997 3Z'
    />
  </svg>
);
