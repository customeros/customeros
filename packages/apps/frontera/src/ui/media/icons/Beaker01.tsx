import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Beaker01 = ({ className, ...props }: IconProps) => (
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
      d='M10 2.00008V6.6605C10 6.87768 10 6.98626 9.9669 7.07264C9.93568 7.15413 9.89435 7.21328 9.82858 7.27063C9.75885 7.33142 9.64593 7.37279 9.42008 7.45553C6.54892 8.50734 4.5 11.2644 4.5 14.5C4.5 18.6421 7.85786 22 12 22C16.1421 22 19.5 18.6421 19.5 14.5C19.5 11.2644 17.4511 8.50734 14.5799 7.45553C14.3541 7.37279 14.2411 7.33142 14.1714 7.27063C14.1056 7.21328 14.0643 7.15413 14.0331 7.07264C14 6.98626 14 6.87768 14 6.6605V2.00008M8.5 2H15.5'
    />
  </svg>
);
