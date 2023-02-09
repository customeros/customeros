import * as React from 'react';
import { SVGProps } from 'react';
const SvgBook = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M19 3.25H6.75a2.43 2.43 0 0 0-2.5 2.35V18a2.85 2.85 0 0 0 2.94 2.75H19a.76.76 0 0 0 .75-.75V4a.76.76 0 0 0-.75-.75Zm-.75 16H7.19A1.35 1.35 0 0 1 5.75 18a1.35 1.35 0 0 1 1.44-1.25h11.06v2.5Zm0-4H7.19a3 3 0 0 0-1.44.37V5.6a.94.94 0 0 1 1-.85h11.5v10.5Z' />
      <path d='M8.75 8.75h6.5a.75.75 0 1 0 0-1.5h-6.5a.75.75 0 0 0 0 1.5ZM8.75 12.25h6.5a.75.75 0 1 0 0-1.5h-6.5a.75.75 0 1 0 0 1.5Z' />
    </g>
  </svg>
);
export default SvgBook;
