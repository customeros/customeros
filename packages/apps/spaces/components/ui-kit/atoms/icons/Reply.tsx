import * as React from 'react';
import { SVGProps } from 'react';
const SvgReply = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M14.55 15.59a.75.75 0 0 1-.55-1.28l3.92-3.89L14 6.53a.75.75 0 0 1 1.06-1.06l4.46 4.42a.75.75 0 0 1 0 1.06l-4.46 4.42a.701.701 0 0 1-.51.22Z' />
      <path d='M5 18.75a.76.76 0 0 1-.75-.75v-7.58A.76.76 0 0 1 5 9.67h14a.75.75 0 1 1 0 1.5H5.75V18a.76.76 0 0 1-.75.75Z' />
    </g>
  </svg>
);
export default SvgReply;
