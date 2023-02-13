import * as React from 'react';
import { SVGProps } from 'react';
const SvgThumbsDown = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12.09 20.75a1.55 1.55 0 0 1-.67-.14 3.16 3.16 0 0 1-1.93-2.91v-2.47h-4a.51.51 0 0 1-.17 0 2.56 2.56 0 0 1-1.58-1 2.29 2.29 0 0 1-.44-1.68l1.1-7.29a2.38 2.38 0 0 1 2.35-2h11.62a2.38 2.38 0 0 1 2.38 2.36v5.64a2.38 2.38 0 0 1-2.38 2.37h-1.95l-2.73 6.09a1.77 1.77 0 0 1-1.6 1.03Zm-6.51-7h4.67a.74.74 0 0 1 .74.75v3.2a1.66 1.66 0 0 0 1 1.53.24.24 0 0 0 .31-.13l2.87-6.39v-8H6.75a.87.87 0 0 0-.87.73l-1.12 7.3a.72.72 0 0 0 .16.56c.162.215.397.364.66.42v.03Zm11.11-1.6h1.68a.87.87 0 0 0 .88-.87V5.61a.87.87 0 0 0-.88-.86h-1.68v7.4Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgThumbsDown;
