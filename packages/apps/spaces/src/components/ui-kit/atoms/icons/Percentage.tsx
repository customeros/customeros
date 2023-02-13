import * as React from 'react';
import { SVGProps } from 'react';
const SvgPercentage = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M7.05 17.7a.739.739 0 0 1-.53-.22.75.75 0 0 1 0-1.06l9.9-9.9a.75.75 0 1 1 1.06 1.06l-9.9 9.9a.74.74 0 0 1-.53.22ZM8.5 10.75a2.25 2.25 0 1 1 0-4.5 2.25 2.25 0 0 1 0 4.5Zm0-3a.75.75 0 1 0 .75.75.76.76 0 0 0-.75-.75ZM15.5 17.75a2.25 2.25 0 1 1 0-4.5 2.25 2.25 0 0 1 0 4.5Zm0-3a.75.75 0 1 0 .75.75.76.76 0 0 0-.75-.75Z' />
    </g>
  </svg>
);
export default SvgPercentage;
