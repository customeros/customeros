import * as React from 'react';
import { SVGProps } from 'react';
const SvgInfo = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 17.75a.76.76 0 0 1-.75-.75v-7a.75.75 0 1 1 1.5 0v7a.76.76 0 0 1-.75.75ZM12 8.25a.76.76 0 0 1-.75-.75V7a.75.75 0 1 1 1.5 0v.5a.76.76 0 0 1-.75.75Z' />
    </g>
  </svg>
);
export default SvgInfo;
