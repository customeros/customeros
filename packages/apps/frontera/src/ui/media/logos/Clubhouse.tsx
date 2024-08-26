import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Clubhouse = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 15'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='#6515DD'
      d='M23.1833 2.66669L12.2166 6.16669V2.70835L0.391602 6.48335V16.5084L10.5916 13.25V16.6917L23.9999 12.4167L20.3666 8.86669L23.1833 2.66669ZM10.5916 11.5334L2.02494 14.2667V7.67502L10.5916 4.94169V11.5334ZM20.9666 11.675L12.2166 14.4667V7.88335L20.2749 5.30835L18.4249 9.20835L20.9666 11.675Z'
    />
    <path
      fill='#6515DD'
      d='M1.74167 17.875C0.783333 17.875 0 18.6583 0 19.6167C0 20.575 0.783333 21.3583 1.74167 21.3583C2.7 21.3583 3.48333 20.575 3.48333 19.6167C3.48333 18.6583 2.70833 17.875 1.74167 17.875Z'
    />
  </svg>
);
