import * as React from 'react';
import { SVGProps } from 'react';
const SvgQuestionCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Zm0 16.5a7.5 7.5 0 1 1 0-15 7.5 7.5 0 0 1 0 15Z' />
      <path d='M9.7 7.7a3.25 3.25 0 0 0-.95 2.3.75.75 0 1 0 1.5 0 1.74 1.74 0 0 1 .51-1.24 1.8 1.8 0 0 1 2.48 0 1.75 1.75 0 0 1-1.24 3 .76.76 0 0 0-.75.75v1a.75.75 0 1 0 1.5 0v-.34a3.19 3.19 0 0 0 1.55-.86 3.26 3.26 0 0 0 0-4.6 3.34 3.34 0 0 0-4.6-.01ZM12 17.5a1.25 1.25 0 1 0 0-2.5 1.25 1.25 0 0 0 0 2.5Z' />
    </g>
  </svg>
);
export default SvgQuestionCircle;
