import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Server03 = ({ className, ...props }: IconProps) => (
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
      d='M22 17.5L21.3083 7.46975C21.1997 5.89504 21.1454 5.10768 20.8041 4.51082C20.5036 3.98533 20.0512 3.56304 19.5062 3.29942C18.8873 3 18.0981 3 16.5196 3H7.48037C5.90191 3 5.11268 3 4.49376 3.29942C3.94884 3.56304 3.49642 3.98533 3.19594 4.51082C2.85464 5.10768 2.80034 5.89503 2.69174 7.46975L2 17.5M22 17.5C22 19.433 20.433 21 18.5 21H5.5C3.567 21 2 19.433 2 17.5M22 17.5C22 15.567 20.433 14 18.5 14H5.5C3.567 14 2 15.567 2 17.5M6 17.5H6.01M12 17.5H18'
    />
  </svg>
);
