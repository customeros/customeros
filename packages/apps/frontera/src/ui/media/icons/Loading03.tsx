import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Loading03 = ({ className, ...props }: IconProps) => (
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
      d='M3.34025 16.9997C2.0881 14.8298 2.72473 12.0644 4.79795 10.6601L4.80018 10.6578C6.46564 9.53235 8.65775 9.57012 10.2843 10.7523L13.7164 13.2477C15.3418 14.4299 17.5339 14.4676 19.2005 13.3421L19.2027 13.3399C21.2771 11.9356 21.9148 9.16792 20.6604 7.00026M17.002 20.6593C14.8321 21.9114 12.0667 21.2748 10.6623 19.2016L10.6601 19.1994C9.53457 17.5339 9.57234 15.3418 10.7545 13.7152L13.2499 10.2832C14.4321 8.6577 14.4699 6.46559 13.3444 4.79901L13.3399 4.79679C11.9356 2.72468 9.16792 2.08582 7.00026 3.3402M19.0705 4.92901C22.9758 8.83436 22.9758 15.1651 19.0705 19.0705C15.1651 22.9758 8.83436 22.9758 4.92901 19.0705C1.02366 15.1651 1.02366 8.83436 4.92901 4.92901C8.83436 1.02366 15.1651 1.02366 19.0705 4.92901Z'
    />
  </svg>
);
