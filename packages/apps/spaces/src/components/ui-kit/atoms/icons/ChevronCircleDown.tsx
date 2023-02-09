import * as React from 'react';
import { SVGProps } from 'react';
const SvgChevronCircleDown = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 21a9 9 0 1 1 0-18 9 9 0 0 1 0 18Zm0-16.5a7.5 7.5 0 1 0 0 15 7.5 7.5 0 0 0 0-15Z' />
      <path d='M12 14.75a.74.74 0 0 1-.53-.22l-4-4a.75.75 0 0 1 1.06-1.06L12 12.94l3.47-3.47a.75.75 0 0 1 1.06 1.06l-4 4a.74.74 0 0 1-.53.22Z' />
    </g>
  </svg>
);
export default SvgChevronCircleDown;
