import * as React from 'react';
import { SVGProps } from 'react';
const SvgKey = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M15 14.75a5.738 5.738 0 0 1-4.07-1.68A5.77 5.77 0 1 1 15 14.75Zm0-10a4.25 4.25 0 1 0 .02 8.5 4.25 4.25 0 0 0-.02-8.5Z' />
      <path d='M4.5 20.25A.75.75 0 0 1 4 19l6.46-6.47a.753.753 0 0 1 1.06 1.07L5 20a.74.74 0 0 1-.5.25Z' />
      <path d='M8 20.75a.739.739 0 0 1-.53-.22l-2-2a.75.75 0 0 1 1.06-1.06l2 2A.75.75 0 0 1 8 20.75ZM10 18.75a.739.739 0 0 1-.53-.22l-2-2a.749.749 0 1 1 1.06-1.06l2 2a.75.75 0 0 1-.53 1.28Z' />
    </g>
  </svg>
);
export default SvgKey;
