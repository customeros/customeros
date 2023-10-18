import * as React from 'react';
import { SVGProps } from 'react';
const SvgCheckWaves = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={36}
    height={36}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect
      x={5.5}
      y={5.5}
      width={25}
      height={25}
      rx={12.5}
      stroke='#079455'
      strokeWidth={1.667}
      opacity={0.3}
    />
    <rect
      x={1.333}
      y={1.333}
      width={33.333}
      height={33.333}
      rx={16.667}
      stroke='#079455'
      strokeWidth={1.667}
      opacity={0.1}
    />
    <g clipPath='url(#checkWaves_svg__a)'>
      <path
        d='m14.25 18 2.5 2.5 5-5m4.583 2.5a8.333 8.333 0 1 1-16.666 0 8.333 8.333 0 0 1 16.666 0Z'
        stroke='#079455'
        strokeWidth={1.667}
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
    <defs>
      <clipPath id='checkWaves_svg__a'>
        <path fill='#fff' transform='translate(8 8)' d='M0 0h20v20H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgCheckWaves;
