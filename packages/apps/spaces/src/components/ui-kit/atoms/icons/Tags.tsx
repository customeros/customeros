import * as React from 'react';
import { SVGProps } from 'react';
const SvgTags = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='m21.07 10.3-6-6a.75.75 0 0 0-.53-.22H3a.76.76 0 0 0-.75.75v7.88c0 .199.08.39.22.53l6 6a2.311 2.311 0 0 0 1.65.68 2.34 2.34 0 0 0 1.64-.68L12 19c.04.082.09.16.15.23a2.33 2.33 0 0 0 3.29 0l5.65-5.66a2.33 2.33 0 0 0-.02-3.27ZM10.7 18.17a.808.808 0 0 1-1.17 0L3.75 12.4V5.58h6.82l5.78 5.78a.83.83 0 0 1 0 1.17l-5.65 5.64Zm9.3-5.65-5.65 5.65a.821.821 0 0 1-1.17 0A.54.54 0 0 0 13 18l4.44-4.45a2.33 2.33 0 0 0 0-3.28l-4.75-4.69h1.54L20 11.36a.822.822 0 0 1 0 1.16Z' />
      <path d='M7 9.75a1.25 1.25 0 1 0 0-2.5 1.25 1.25 0 0 0 0 2.5Z' />
    </g>
  </svg>
);
export default SvgTags;
