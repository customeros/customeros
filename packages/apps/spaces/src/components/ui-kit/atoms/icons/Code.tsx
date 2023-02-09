import * as React from 'react';
import { SVGProps } from 'react';
const SvgCode = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M9.53 6.47a.75.75 0 0 0-1.06 0l-5 5a.75.75 0 0 0 0 1.06l5 5a.75.75 0 0 0 1.06-1.06L5.06 12l4.47-4.47a.75.75 0 0 0 0-1.06ZM20.53 11.47l-5-5a.75.75 0 0 0-1.06 1.06L18.94 12l-4.47 4.47a.75.75 0 0 0 1.06 1.06l5-5a.75.75 0 0 0 0-1.06Z' />
    </g>
  </svg>
);
export default SvgCode;
