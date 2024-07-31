import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const WifiOff = ({ className, ...props }: IconProps) => (
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
      d='M15.3119 10C16.6802 10.4263 17.9624 11.1191 19.08 12.05M22.5799 8.49997C19.6575 5.92394 15.8956 4.50262 11.9999 4.50262C11.3949 4.50262 10.7931 4.5369 10.1972 4.60447M8.52979 15.61C9.54499 14.8888 10.7595 14.5013 12.0048 14.5013C13.2501 14.5013 14.4646 14.8888 15.4798 15.61M12 19.5H12.01M1.19336 8.70076C2.52697 7.47869 4.06839 6.47975 5.75851 5.76306M4.73193 12.243C6.12934 11.012 7.84172 10.1302 9.73265 9.73393M15.6983 15.7751C14.6792 14.9763 13.3952 14.5 11.9999 14.5C10.5835 14.5 9.28172 14.9908 8.25537 15.8116M3 3L21 21'
    />
  </svg>
);
