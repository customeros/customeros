import * as React from 'react';
import { SVGProps } from 'react';
const SvgInfoCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 16.75a.76.76 0 0 1-.75-.75v-5a.75.75 0 1 1 1.5 0v5a.76.76 0 0 1-.75.75ZM12 9.25a.76.76 0 0 1-.75-.75V8a.75.75 0 1 1 1.5 0v.5a.76.76 0 0 1-.75.75Z' />
      <path d='M12 21a9 9 0 1 1 0-18 9 9 0 0 1 0 18Zm0-16.5a7.5 7.5 0 1 0 0 15 7.5 7.5 0 0 0 0-15Z' />
    </g>
  </svg>
);
export default SvgInfoCircle;
