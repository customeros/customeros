import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Seeding = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='2'
      stroke='currentColor'
      d='M11.875 14.25V13.25H10.875C7.07846 13.25 4 10.1715 4 6.375C4 6.30502 4.05502 6.25 4.125 6.25H5.25C9.04654 6.25 12.125 9.32846 12.125 13.125V14.25V19.875C12.125 19.945 12.07 20 12 20C11.93 20 11.875 19.945 11.875 19.875V14.25ZM20 4.125C20 7.28352 17.8684 9.94659 14.9673 10.7499C14.6237 9.33931 13.9819 8.04611 13.1126 6.94205C14.3556 5.16265 16.4179 4 18.75 4H19.875C19.945 4 20 4.05502 20 4.125Z'
    />
  </svg>
);
