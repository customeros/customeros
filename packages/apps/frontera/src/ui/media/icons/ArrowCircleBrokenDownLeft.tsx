import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ArrowCircleBrokenDownLeft = ({
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
      d='M9.00023 9.00001V15M9.00023 15H15.0002M9.00023 15L19 4.99995M21.6606 9.41045C22.5515 12.7466 21.6884 16.4538 19.0711 19.0711C15.1658 22.9763 8.83418 22.9763 4.92893 19.0711C1.02369 15.1658 1.02369 8.83418 4.92893 4.92893C7.54623 2.31164 11.2534 1.44846 14.5896 2.33938'
    />
  </svg>
);
