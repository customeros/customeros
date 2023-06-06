import * as React from 'react';
import { SVGProps } from 'react';
const SvgDataProvider = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g
      clipPath='url(#data-provider_svg__a)'
      stroke='#000'
      strokeMiterlimit={10}
    >
      <path d='M18.7 9.11h-1V7.2A5.75 5.75 0 0 0 12 1.44h-.21A5.75 5.75 0 0 0 7 4.41l-.13.26a6.5 6.5 0 0 0-.62 2.53A4.78 4.78 0 0 0 1.44 12a1.27 1.27 0 0 0 0 .2 4.8 4.8 0 0 0 4.78 4.6H18.7a3.83 3.83 0 0 0 2.72-1.13 3.9 3.9 0 0 0 1.11-2.45 2.255 2.255 0 0 0 0-.26 3.84 3.84 0 0 0-3.83-3.85ZM11.99 16.79v6.71' />
      <path d='M4.32 22.54h1.91a1.92 1.92 0 0 0 1.92-1.92v-3.83M15.83 16.79v3.83a1.91 1.91 0 0 0 1.91 1.92h1.92' />
    </g>
    <defs>
      <clipPath id='data-provider_svg__a'>
        <path fill='#fff' d='M0 0h24v24H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgDataProvider;
