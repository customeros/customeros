import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Telegram = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='url(#paint0_linear_1027_1570)'
      d='M12 24C18.6274 24 24 18.6274 24 12C24 5.37258 18.6274 0 12 0C5.37258 0 0 5.37258 0 12C0 18.6274 5.37258 24 12 24Z'
    />
    <path
      fill='white'
      fill-rule='evenodd'
      clip-rule='evenodd'
      d='M5.43201 11.8733C8.93026 10.3492 11.263 9.3444 12.4301 8.85893C15.7627 7.47282 16.4551 7.23203 16.9065 7.22408C17.0058 7.22234 17.2277 7.24694 17.3715 7.3636C17.4929 7.46211 17.5263 7.59518 17.5423 7.68857C17.5583 7.78197 17.5782 7.99473 17.5623 8.16097C17.3817 10.0585 16.6003 14.6631 16.2028 16.7884C16.0346 17.6876 15.7034 17.9891 15.3827 18.0186C14.6858 18.0828 14.1567 17.5581 13.4817 17.1157C12.4256 16.4233 11.8289 15.9924 10.8037 15.3168C9.61896 14.536 10.387 14.1069 11.0622 13.4056C11.2389 13.2221 14.3093 10.4294 14.3687 10.176C14.3762 10.1443 14.3831 10.0262 14.3129 9.96385C14.2427 9.90148 14.1392 9.92281 14.0644 9.93977C13.9585 9.96381 12.2713 11.079 9.00276 13.2853C8.52385 13.6142 8.09007 13.7744 7.70141 13.766C7.27295 13.7568 6.44876 13.5238 5.83606 13.3246C5.08456 13.0803 4.48728 12.9512 4.53929 12.5363C4.56638 12.3202 4.86395 12.0992 5.43201 11.8733Z'
    />
    <defs>
      <linearGradient
        y1='0'
        x1='12'
        x2='12'
        y2='23.822'
        id='paint0_linear_1027_1570'
        gradientUnits='userSpaceOnUse'
      >
        <stop stop-color='#2AABEE' />
        <stop offset='1' stop-color='#229ED9' />
      </linearGradient>
    </defs>
  </svg>
);
