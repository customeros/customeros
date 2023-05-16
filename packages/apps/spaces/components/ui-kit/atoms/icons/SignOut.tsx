import * as React from 'react';
import { SVGProps } from 'react';
const SvgSignOut = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M9 20.75H6a2.64 2.64 0 0 1-2.75-2.53V5.78A2.64 2.64 0 0 1 6 3.25h3a.75.75 0 1 1 0 1.5H6a1.16 1.16 0 0 0-1.25 1v12.47a1.16 1.16 0 0 0 1.25 1h3a.75.75 0 1 1 0 1.5v.03zm7-4a.74.74 0 0 1-.53-.22.75.75 0 0 1 0-1.06L18.94 12l-3.47-3.47a.75.75 0 0 1-.019-1.079.75.75 0 0 1 1.079.019l4 4a.75.75 0 0 1 0 1.06l-4 4a.74.74 0 0 1-.53.22z' />
    <path d='M20 12.75H9a.75.75 0 0 1-.75-.75.75.75 0 0 1 .75-.75h11a.75.75 0 0 1 .75.75.75.75 0 0 1-.75.75z' />
  </svg>
);
export default SvgSignOut;
