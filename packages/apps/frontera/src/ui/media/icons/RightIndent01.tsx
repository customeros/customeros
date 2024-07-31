import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const RightIndent01 = ({ className, ...props }: IconProps) => (
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
      d='M21 4H3M21 20H3M12 9.25H3M12 14.75H3M19.72 8.46L15.8533 11.36C15.5638 11.5771 15.4191 11.6857 15.3673 11.8188C15.322 11.9353 15.322 12.0647 15.3673 12.1812C15.4191 12.3143 15.5638 12.4229 15.8533 12.64L19.72 15.54C20.132 15.849 20.338 16.0035 20.5104 15.9999C20.6604 15.9968 20.8012 15.9264 20.8937 15.8082C21 15.6725 21 15.415 21 14.9V9.1C21 8.58503 21 8.32754 20.8937 8.19175C20.8012 8.07358 20.6604 8.00323 20.5104 8.00011C20.338 7.99652 20.132 8.15102 19.72 8.46Z'
    />
  </svg>
);
