import * as React from 'react';
import { SVGProps } from 'react';
const SvgRefresh = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 20.75a7.25 7.25 0 1 1 0-14.5h2.5a.75.75 0 1 1 0 1.5H12a5.75 5.75 0 1 0 5.75 5.75.75.75 0 1 1 1.5 0A7.26 7.26 0 0 1 12 20.75Z' />
      <path d='M12 10.75a.741.741 0 0 1-.53-.22.75.75 0 0 1 0-1.06L13.94 7l-2.47-2.47a.748.748 0 0 1 .23-1.244.75.75 0 0 1 .83.184l3 3a.75.75 0 0 1 0 1.06l-3 3a.738.738 0 0 1-.53.22Z' />
    </g>
  </svg>
);
export default SvgRefresh;
