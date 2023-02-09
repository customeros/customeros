import * as React from 'react';
import { SVGProps } from 'react';
const SvgAlignLeft = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M13.93 11h-10a.75.75 0 1 1 0-1.5h10a.75.75 0 1 1 0 1.5ZM20.07 7h-16a.75.75 0 0 1 0-1.5h16a.75.75 0 1 1 0 1.5ZM20.07 15h-16a.75.75 0 1 1 0-1.5h16a.75.75 0 1 1 0 1.5ZM13.93 19h-10a.75.75 0 1 1 0-1.5h10a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgAlignLeft;
