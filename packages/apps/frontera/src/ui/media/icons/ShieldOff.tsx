import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ShieldOff = ({ className, ...props }: IconProps) => (
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
      d='M8.80868 3.19661L11.4382 2.21054C11.6461 2.13258 11.75 2.0936 11.857 2.07815C11.9518 2.06444 12.0482 2.06444 12.143 2.07815C12.25 2.0936 12.3539 2.13258 12.5618 2.21054L17.9236 4.22122C18.6722 4.50195 19.0465 4.64231 19.3223 4.8854C19.566 5.10015 19.7537 5.37101 19.8692 5.67457C20 6.0182 20 6.41796 20 7.21747V11.9999C20 12.7202 19.8847 13.4117 19.6818 14.0708M17.6219 17.6255C15.9362 19.5859 13.7879 20.9788 12.6996 21.6138C12.4772 21.7436 12.366 21.8085 12.2098 21.842C12.0885 21.868 11.9116 21.868 11.7904 21.842C11.6341 21.8085 11.5234 21.7439 11.302 21.6147C9.35396 20.4783 4 16.9083 4 11.9999V5.77618C4 5.30923 4.28966 4.89125 4.72688 4.72729M3 2.99987L21 20.9999'
    />
  </svg>
);
