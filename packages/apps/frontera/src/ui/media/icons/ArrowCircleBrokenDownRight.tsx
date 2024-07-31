import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ArrowCircleBrokenDownRight = ({
  className,
  ...props
}: IconProps) => (
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
      d='M9.41045 2.33938C12.7466 1.44846 16.4538 2.31164 19.0711 4.92893C22.9763 8.83418 22.9763 15.1658 19.0711 19.0711C15.1658 22.9763 8.83418 22.9763 4.92893 19.0711C2.31164 16.4538 1.44846 12.7466 2.33938 9.41045M15.0001 9.00001V15M15.0001 15H9.00014M15.0001 15L4.99995 4.99995'
    />
  </svg>
);
