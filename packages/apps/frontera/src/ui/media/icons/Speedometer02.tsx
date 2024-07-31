import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Speedometer02 = ({ className, ...props }: IconProps) => (
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
      d='M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12M22 12C22 6.47715 17.5228 2 12 2M22 12H19.5M2 12C2 6.47715 6.47715 2 12 2M2 12H4.5M12 2V4.5M19.0784 5L13.4999 10.5M19.0784 19.0784L18.8745 18.8745C18.1827 18.1827 17.8368 17.8368 17.4331 17.5894C17.0753 17.3701 16.6851 17.2085 16.2769 17.1105C15.8166 17 15.3274 17 14.349 17L9.65096 17C8.6726 17 8.18342 17 7.72307 17.1106C7.31493 17.2086 6.92475 17.3702 6.56686 17.5895C6.1632 17.8369 5.8173 18.1828 5.12549 18.8746L4.92163 19.0784M4.92163 5L6.65808 6.73645M14 12C14 13.1046 13.1046 14 12 14C10.8954 14 10 13.1046 10 12C10 10.8954 10.8954 10 12 10C13.1046 10 14 10.8954 14 12Z'
    />
  </svg>
);
