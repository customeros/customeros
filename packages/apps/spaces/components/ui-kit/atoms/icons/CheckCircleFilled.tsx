import * as React from 'react';
import { SVGProps } from 'react';
const SvgCheckCircleFilled = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 6 6'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      fillRule='evenodd'
      clipRule='evenodd'
      d='M3 6a3 3 0 0 0 2.796-4.089L3 4.707 1.146 2.854l.708-.708L3 3.293 5.263 1.03A3 3 0 1 0 3 6Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCheckCircleFilled;
