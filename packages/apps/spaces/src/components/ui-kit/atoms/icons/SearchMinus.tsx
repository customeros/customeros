import * as React from 'react';
import { SVGProps } from 'react';
const SvgSearchMinus = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M10.77 18.3a7.53 7.53 0 1 1 7.53-7.53 7.54 7.54 0 0 1-7.53 7.53Zm0-13.55a6 6 0 1 0 0 12 6 6 0 0 0 0-12Z' />
      <path d='M20 20.75a.741.741 0 0 1-.53-.22l-4.13-4.13a.75.75 0 0 1 1.06-1.06l4.13 4.13a.75.75 0 0 1-.53 1.28ZM13.25 11.5h-5a.75.75 0 1 1 0-1.5h5a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgSearchMinus;
