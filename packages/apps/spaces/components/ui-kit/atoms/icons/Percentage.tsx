import * as React from 'react';
import { SVGProps } from 'react';
const SvgPercentage = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M7.05 17.7a.74.74 0 0 1-.53-.22.75.75 0 0 1 0-1.06l9.9-9.9a.75.75 0 0 1 1.078-.019.75.75 0 0 1-.019 1.079l-9.9 9.9a.74.74 0 0 1-.53.22zm1.45-6.95a2.25 2.25 0 0 1-1.591-3.841A2.25 2.25 0 0 1 10.75 8.5a2.25 2.25 0 0 1-2.25 2.25zm0-3a.75.75 0 0 0-.53 1.28.75.75 0 0 0 1.28-.53.76.76 0 0 0-.75-.75zm7 10a2.25 2.25 0 0 1-1.591-3.841A2.25 2.25 0 0 1 17.75 15.5a2.25 2.25 0 0 1-2.25 2.25zm0-3a.75.75 0 0 0-.53 1.28.75.75 0 0 0 1.28-.53.76.76 0 0 0-.75-.75z' />
  </svg>
);
export default SvgPercentage;
