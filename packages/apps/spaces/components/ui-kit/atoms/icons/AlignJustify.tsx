import * as React from 'react';
import { SVGProps } from 'react';
const SvgAlignJustify = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M20 10.75H4a.75.75 0 1 1 0-1.5h16a.75.75 0 1 1 0 1.5ZM20 6.75H4a.75.75 0 0 1 0-1.5h16a.75.75 0 1 1 0 1.5ZM20 14.75H4a.75.75 0 1 1 0-1.5h16a.75.75 0 1 1 0 1.5ZM20 18.75H4a.75.75 0 1 1 0-1.5h16a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgAlignJustify;
