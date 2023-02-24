import * as React from 'react';
import { SVGProps } from 'react';
const SvgCheck = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M9 18.25a.74.74 0 0 1-.53-.25l-5-5a.75.75 0 1 1 1.06-1L9 16.44 19.47 6a.75.75 0 0 1 1.06 1l-11 11a.74.74 0 0 1-.53.25Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCheck;
