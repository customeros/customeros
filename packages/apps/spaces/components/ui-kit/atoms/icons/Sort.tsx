import * as React from 'react';
import { SVGProps } from 'react';
const SvgSort = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M18 10.75H6a.74.74 0 0 1-.69-.46.75.75 0 0 1 .16-.82l6-6a.75.75 0 0 1 1.06 0l6 6a.75.75 0 0 1 .16.82.74.74 0 0 1-.69.46zM7.81 9.25h8.38L12 5.06 7.81 9.25zM12 20.75a.74.74 0 0 1-.53-.22l-6-6a.75.75 0 0 1-.16-.82.74.74 0 0 1 .69-.46h12a.74.74 0 0 1 .69.46.75.75 0 0 1-.16.82l-6 6a.74.74 0 0 1-.53.22zm-4.19-6L12 18.94l4.19-4.19H7.81z' />
  </svg>
);
export default SvgSort;
