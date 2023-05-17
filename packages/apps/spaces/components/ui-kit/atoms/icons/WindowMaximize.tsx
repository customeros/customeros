import * as React from 'react';
import { SVGProps } from 'react';
const SvgWindowMaximize = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M18 20.75h-6a.75.75 0 1 1 0-1.5h6A1.25 1.25 0 0 0 19.25 18V6A1.25 1.25 0 0 0 18 4.75H6A1.25 1.25 0 0 0 4.75 6v6a.75.75 0 1 1-1.5 0V6A2.75 2.75 0 0 1 6 3.25h12A2.75 2.75 0 0 1 20.75 6v12A2.75 2.75 0 0 1 18 20.75zm-2-8a.76.76 0 0 1-.75-.75V8.75H12a.75.75 0 1 1 0-1.5h4a.76.76 0 0 1 .75.75v4a.76.76 0 0 1-.75.75z' />
    <path d='M11.5 13.25A.74.74 0 0 1 11 13a.75.75 0 0 1 0-1l4.5-4.5a.75.75 0 0 1 1.06 1.06L12 13a.74.74 0 0 1-.5.25zM8 20.75H5A1.76 1.76 0 0 1 3.25 19v-3A1.76 1.76 0 0 1 5 14.25h3A1.76 1.76 0 0 1 9.75 16v3A1.76 1.76 0 0 1 8 20.75zm-3-5a.25.25 0 0 0-.25.25v3a.25.25 0 0 0 .25.25h3a.25.25 0 0 0 .25-.25v-3a.25.25 0 0 0-.25-.25H5z' />
  </svg>
);
export default SvgWindowMaximize;
