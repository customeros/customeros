import * as React from 'react';
import { SVGProps } from 'react';
const SvgClockCheck = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g clipPath='url(#clock-check_svg__a)'>
      <path
        d='M9.667 12.667 11 14l3-3m.657-2.633a6.667 6.667 0 1 0-6.831 6.298M8 4v4l2.492 1.246'
        stroke='currentColor'
        strokeWidth={1.333}
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
    <defs>
      <clipPath id='clock-check_svg__a'>
        <path fill='#fff' d='M0 0h16v16H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgClockCheck;
