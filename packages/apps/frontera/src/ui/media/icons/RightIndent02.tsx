import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const RightIndent02 = ({ className, ...props }: IconProps) => (
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
      d='M12 9.24995H3M12 3.99995L3 3.99995M21 14.75H3M21 20H3M19.72 2.95995L15.8533 5.85995C15.5638 6.07707 15.4191 6.18563 15.3673 6.31872C15.322 6.43529 15.322 6.56462 15.3673 6.68119C15.4191 6.81427 15.5638 6.92283 15.8533 7.13995L19.72 10.04C20.132 10.3489 20.338 10.5034 20.5104 10.4998C20.6604 10.4967 20.8012 10.4264 20.8937 10.3082C21 10.1724 21 9.91493 21 9.39995V3.59995C21 3.08498 21 2.82749 20.8937 2.6917C20.8012 2.57354 20.6604 2.50318 20.5104 2.50006C20.338 2.49648 20.132 2.65097 19.72 2.95995Z'
    />
  </svg>
);
