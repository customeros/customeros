import * as React from 'react';
import { SVGProps } from 'react';
const SvgCelandar = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17 4.25h-1.25V3a.75.75 0 1 0-1.5 0v1.25h-4.5V3a.75.75 0 0 0-1.5 0v1.25H7A2.75 2.75 0 0 0 4.25 7v11A2.75 2.75 0 0 0 7 20.75h10A2.75 2.75 0 0 0 19.75 18V7A2.75 2.75 0 0 0 17 4.25zM7 5.75h1.25V7a.75.75 0 0 0 1.5 0V5.75h4.5V7a.75.75 0 1 0 1.5 0V5.75H17A1.25 1.25 0 0 1 18.25 7v2.75H5.75V7A1.25 1.25 0 0 1 7 5.75zm10 13.5H7A1.25 1.25 0 0 1 5.75 18v-6.75h12.5V18A1.25 1.25 0 0 1 17 19.25z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCelandar;
