import * as React from 'react';
import { SVGProps } from 'react';
const SvgTablet = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M18 3.25H6A1.76 1.76 0 0 0 4.25 5v14A1.76 1.76 0 0 0 6 20.75h12A1.76 1.76 0 0 0 19.75 19V5A1.76 1.76 0 0 0 18 3.25ZM18.25 19a.25.25 0 0 1-.25.25H6a.25.25 0 0 1-.25-.25V5A.25.25 0 0 1 6 4.75h12a.25.25 0 0 1 .25.25v14Z' />
      <path d='M12 14.5a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3Z' />
    </g>
  </svg>
);
export default SvgTablet;
