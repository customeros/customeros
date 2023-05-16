import * as React from 'react';
import { SVGProps } from 'react';
const SvgTag = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='m19.32 10.72-6.48-6.48a.75.75 0 0 0-.53-.24H4.77a.76.76 0 0 0-.77.77v7.54a.79.79 0 0 0 .22.54l6.48 6.47a2.26 2.26 0 0 0 3.19 0l5.41-5.41a2.26 2.26 0 0 0 .02-3.19zm-1.06 2.13-5.41 5.41a.77.77 0 0 1-1.07 0L5.52 12V5.52H12l6.26 6.26a.77.77 0 0 1 0 1.07zM8.5 9.75a1.25 1.25 0 1 0 0-2.5 1.25 1.25 0 1 0 0 2.5z' />
  </svg>
);
export default SvgTag;
