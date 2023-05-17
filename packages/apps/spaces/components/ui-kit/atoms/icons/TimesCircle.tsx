import * as React from 'react';
import { SVGProps } from 'react';
const SvgTimesCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12 21A9 9 0 0 1 5.636 5.636 9 9 0 0 1 21 12a9 9 0 0 1-9 9zm0-16.5a7.5 7.5 0 0 0-5.303 12.803A7.5 7.5 0 0 0 19.5 12 7.5 7.5 0 0 0 12 4.5zM9 15.75a.74.74 0 0 1-.53-.22.75.75 0 0 1 0-1.06l6-6a.75.75 0 0 1 1.06 1.06l-6 6a.74.74 0 0 1-.53.22z' />
    <path d='M15 15.75a.74.74 0 0 1-.53-.22l-6-6a.75.75 0 0 1 .018-1.042.75.75 0 0 1 1.042-.018l6 6a.75.75 0 0 1 0 1.06.74.74 0 0 1-.53.22z' />
  </svg>
);
export default SvgTimesCircle;
