import * as React from 'react';
import { SVGProps } from 'react';
const SvgPhoneSlashed = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    fill='none'
    viewBox='0 0 24 24'
    {...props}
  >
    <g
      stroke='#fff'
      strokeLinecap='round'
      strokeLinejoin='round'
      strokeWidth={2}
      clipPath='url(#phone-slashed_svg__a)'
    >
      <path d='M5.2 12.8a19.8 19.8 0 0 1-3-8.6 2 2 0 0 1 2-2.2h3a2 2 0 0 1 2 1.7c0 1 .3 2 .6 2.8a2 2 0 0 1-.4 2.1L8 10m2.6 3.4c1 1 2.1 1.9 3.4 2.6l1.3-1.3a2 2 0 0 1 2-.4c1 .3 2 .6 2.9.7a2 2 0 0 1 1.7 2v3a2 2 0 0 1-2.2 2 19.8 19.8 0 0 1-12-5.8l2.9-2.8ZM23 1 1 23' />
    </g>
    <defs>
      <clipPath id='phone-slashed_svg__a'>
        <path fill='#fff' d='M0 0h24v24H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgPhoneSlashed;
