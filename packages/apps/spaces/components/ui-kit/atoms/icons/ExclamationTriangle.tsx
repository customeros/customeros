import * as React from 'react';
import { SVGProps } from 'react';
const SvgExclamationTriangle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M20 18.75H4a.76.76 0 0 1-.65-.37.77.77 0 0 1 0-.75l8-14a.78.78 0 0 1 1.3 0l8 14a.77.77 0 0 1 0 .75.76.76 0 0 1-.65.37Zm-14.71-1.5h13.42L12 5.51 5.29 17.25Z'
      fill='currentColor'
    />
    <path
      d='M12 13.25a.76.76 0 0 1-.75-.75V9a.75.75 0 1 1 1.5 0v3.5a.76.76 0 0 1-.75.75ZM12 16.25a.76.76 0 0 1-.75-.75V15a.75.75 0 1 1 1.5 0v.5a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgExclamationTriangle;
