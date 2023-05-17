import * as React from 'react';
import { SVGProps } from 'react';
const SvgExclamationTriangle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M20 18.75H4a.76.76 0 0 1-.65-.37.77.77 0 0 1 0-.75l8-14a.78.78 0 0 1 1.3 0l8 14a.77.77 0 0 1 0 .75.76.76 0 0 1-.65.37zm-14.71-1.5h13.42L12 5.51 5.29 17.25zm6.71-4a.76.76 0 0 1-.75-.75V9a.75.75 0 1 1 1.5 0v3.5a.76.76 0 0 1-.75.75zm0 3a.76.76 0 0 1-.75-.75V15a.75.75 0 1 1 1.5 0v.5a.76.76 0 0 1-.75.75z' />
  </svg>
);
export default SvgExclamationTriangle;
