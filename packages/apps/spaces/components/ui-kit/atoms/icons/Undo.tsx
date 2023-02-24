import * as React from 'react';
import { SVGProps } from 'react';
const SvgUndo = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 20.75a7.26 7.26 0 0 1-7.25-7.25.75.75 0 1 1 1.5 0A5.75 5.75 0 1 0 12 7.75H9.5a.75.75 0 0 1 0-1.5H12a7.25 7.25 0 1 1 0 14.5Z' />
      <path d='M12 10.75a.74.74 0 0 1-.53-.22l-3-3a.75.75 0 0 1 0-1.06l3-3a.749.749 0 1 1 1.06 1.06L10.06 7l2.47 2.47a.75.75 0 0 1-.53 1.28Z' />
    </g>
  </svg>
);
export default SvgUndo;
