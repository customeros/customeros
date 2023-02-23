import * as React from 'react';
import { SVGProps } from 'react';
const SvgCreditCard = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M19 5.25H5A1.76 1.76 0 0 0 3.25 7v10A1.76 1.76 0 0 0 5 18.75h14A1.76 1.76 0 0 0 20.75 17V7A1.76 1.76 0 0 0 19 5.25ZM5 6.75h14a.25.25 0 0 1 .25.25v2.25H4.75V7A.25.25 0 0 1 5 6.75Zm14 10.5H5a.25.25 0 0 1-.25-.25v-6.25h14.5V17a.25.25 0 0 1-.25.25Z' />
      <path d='M9 13H7a1 1 0 0 0 0 2h2a1 1 0 0 0 0-2Z' />
    </g>
  </svg>
);
export default SvgCreditCard;
