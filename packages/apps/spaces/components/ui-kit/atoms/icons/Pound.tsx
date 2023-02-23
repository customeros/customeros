import * as React from 'react';
import { SVGProps } from 'react';
const SvgPound = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17.54 19.25H8.1l2-2.82a.76.76 0 0 0 .14-.43v-3.25h3.25a.75.75 0 1 0 0-1.5h-3.28V8a3.09 3.09 0 0 1 3.17-3.25 3.14 3.14 0 0 1 3.33 3.41v1a.75.75 0 1 0 1.5 0v-1a4.621 4.621 0 0 0-4.83-4.91A4.57 4.57 0 0 0 8.71 8v3.25H6.46a.75.75 0 1 0 0 1.5h2.25v3l-2.66 3.82a.76.76 0 0 0 0 .78.74.74 0 0 0 .66.4h10.83a.75.75 0 1 0 0-1.5Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgPound;
