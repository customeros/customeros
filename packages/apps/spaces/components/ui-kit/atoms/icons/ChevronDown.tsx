import * as React from 'react';
import { SVGProps } from 'react';
const SvgChevronDown = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 15.25a.74.74 0 0 1-.53-.22l-5-5A.75.75 0 0 1 7.53 9L12 13.44 16.47 9a.75.75 0 0 1 1.06 1l-5 5a.74.74 0 0 1-.53.25Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgChevronDown;
