import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cursor01 = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='white'
      fillRule='evenodd'
      clipRule='evenodd'
      d='M12.9878 22.0522L15.1245 20.9018L16.9768 19.9396L14.0315 14.416H18.9999L5.94922 1.33203V19.6999L9.75238 16.0056L12.9878 22.0522Z'
    />
    <path
      fillRule='evenodd'
      clipRule='evenodd'
      fill='currentColor'
      d='M13.3124 20.3634L15.3367 19.2841L12.154 13.3179H16.2875L7.0835 4.09326V16.9249L9.98519 14.1242L13.3124 20.3634Z'
    />
  </svg>
);
