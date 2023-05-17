import * as React from 'react';
import { SVGProps } from 'react';
const SvgSlidersH = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M19 8.25h-7a.75.75 0 1 1 0-1.5h7a.75.75 0 1 1 0 1.5zm-9 0H5a.75.75 0 1 1 0-1.5h5a.75.75 0 1 1 0 1.5z' />
    <path d='M10 9.75A.76.76 0 0 1 9.25 9V6a.75.75 0 1 1 1.5 0v3a.76.76 0 0 1-.75.75zm9 7.5h-7a.75.75 0 1 1 0-1.5h7a.75.75 0 1 1 0 1.5zm-9 0H5a.75.75 0 1 1 0-1.5h5a.75.75 0 1 1 0 1.5z' />
    <path d='M10 18.75a.76.76 0 0 1-.75-.75v-3a.75.75 0 1 1 1.5 0v3a.76.76 0 0 1-.75.75zm9-6h-3a.75.75 0 1 1 0-1.5h3a.75.75 0 1 1 0 1.5zm-5 0H5a.75.75 0 1 1 0-1.5h9a.75.75 0 1 1 0 1.5z' />
    <path d='M14 14.25a.76.76 0 0 1-.75-.75v-3a.75.75 0 0 1 .75-.75.75.75 0 0 1 .75.75v3a.76.76 0 0 1-.75.75z' />
  </svg>
);
export default SvgSlidersH;
