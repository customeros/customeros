import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const InvoiceStart = ({ className, ...props }: IconProps) => (
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
      d='M15.2 3H8.8C7.11984 3 6.27976 3 5.63803 3.32698C5.07354 3.6146 4.6146 4.07354 4.32698 4.63803C4 5.27976 4 6.11984 4 7.8V21L6.75 19L9.25 21L12 19L14.75 21L17.25 19L20 21V7.8C20 6.11984 20 5.27976 19.673 4.63803C19.3854 4.07354 18.9265 3.6146 18.362 3.32698C17.7202 3 16.8802 3 15.2 3Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M9.08331 8.0788C9.08331 7.67415 9.08331 7.47182 9.16768 7.36029C9.24119 7.26313 9.35353 7.20301 9.47515 7.19575C9.61474 7.18741 9.78309 7.29964 10.1198 7.5241L14.5013 10.4451C14.7795 10.6306 14.9186 10.7233 14.967 10.8402C15.0094 10.9424 15.0094 11.0572 14.967 11.1594C14.9186 11.2763 14.7795 11.369 14.5013 11.5545L10.1198 14.4755C9.78309 14.6999 9.61474 14.8122 9.47515 14.8038C9.35353 14.7966 9.24119 14.7364 9.16768 14.6393C9.08331 14.5278 9.08331 14.3254 9.08331 13.9208V8.0788Z'
    />
  </svg>
);
