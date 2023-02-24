import * as React from 'react';
import { SVGProps } from 'react';
const SvgQuestion = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M16.07 4.93A5.75 5.75 0 0 0 6.25 9a.75.75 0 1 0 1.5 0A4.26 4.26 0 1 1 12 13.25a.76.76 0 0 0-.75.75v2a.75.75 0 1 0 1.5 0v-1.3a5.76 5.76 0 0 0 3.32-9.77ZM12 20.75a1.25 1.25 0 1 0 0-2.5 1.25 1.25 0 0 0 0 2.5Z' />
    </g>
  </svg>
);
export default SvgQuestion;
