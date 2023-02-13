import * as React from 'react';
import { SVGProps } from 'react';
const SvgTimesCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 21a9 9 0 1 1 0-18 9 9 0 0 1 0 18Zm0-16.5a7.5 7.5 0 1 0 0 15 7.5 7.5 0 0 0 0-15Z' />
      <path d='M9 15.75a.739.739 0 0 1-.53-.22.75.75 0 0 1 0-1.06l6-6a.75.75 0 0 1 1.06 1.06l-6 6a.74.74 0 0 1-.53.22Z' />
      <path d='M15 15.75a.74.74 0 0 1-.53-.22l-6-6a.75.75 0 0 1 1.06-1.06l6 6a.75.75 0 0 1-.53 1.28Z' />
    </g>
  </svg>
);
export default SvgTimesCircle;
