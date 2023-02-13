import * as React from 'react';
import { SVGProps } from 'react';
const SvgArrowUp = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M18 11.75a.74.74 0 0 1-.53-.22L12 6.06l-5.47 5.47a.75.75 0 0 1-1.06-1.06l6-6a.75.75 0 0 1 1.06 0l6 6a.75.75 0 0 1-.53 1.28Z' />
      <path d='M12 19.75a.76.76 0 0 1-.75-.75V5a.75.75 0 1 1 1.5 0v14a.76.76 0 0 1-.75.75Z' />
    </g>
  </svg>
);
export default SvgArrowUp;
