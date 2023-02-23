import * as React from 'react';
import { SVGProps } from 'react';
const SvgSortDown = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 16.25a.74.74 0 0 1-.53-.22l-7-7A.75.75 0 0 1 5 7.75h14A.75.75 0 0 1 19.53 9l-7 7a.74.74 0 0 1-.53.25Zm-5.19-7L12 14.44l5.19-5.19H6.81Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgSortDown;
