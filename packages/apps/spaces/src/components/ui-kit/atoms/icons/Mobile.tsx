import * as React from 'react';
import { SVGProps } from 'react';
const SvgMobile = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M16 3.25H8A1.76 1.76 0 0 0 6.25 5v14A1.76 1.76 0 0 0 8 20.75h8A1.76 1.76 0 0 0 17.75 19V5A1.76 1.76 0 0 0 16 3.25ZM16.25 19a.25.25 0 0 1-.25.25H8a.25.25 0 0 1-.25-.25V5A.25.25 0 0 1 8 4.75h8a.25.25 0 0 1 .25.25v14Z' />
      <path d='M12 14.5a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3Z' />
    </g>
  </svg>
);
export default SvgMobile;
