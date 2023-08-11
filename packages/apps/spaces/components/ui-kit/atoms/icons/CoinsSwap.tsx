import * as React from 'react';
import { SVGProps } from 'react';
const SvgCoinsSwap = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g clipPath='url(#coins-swap_svg__a)'>
      <path
        d='m4 4 1.333-1.333m0 0L4 1.334m1.333 1.333H4a2.667 2.667 0 0 0-2.667 2.666M12 12l-1.333 1.333m0 0L12 14.667m-1.333-1.334H12a2.667 2.667 0 0 0 2.667-2.666M8.945 8.945a4 4 0 1 0-1.89-1.89m2.278 3.612a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z'
        stroke='currentColor'
        strokeWidth={1.333}
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
    <defs>
      <clipPath id='coins-swap_svg__a'>
        <path fill='#fff' d='M0 0h16v16H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgCoinsSwap;
