import * as React from 'react';
import { SVGProps } from 'react';
const SvgPause = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M9 19.75a.76.76 0 0 1-.75-.75V5a.75.75 0 0 1 1.5 0v14a.76.76 0 0 1-.75.75ZM15 19.75a.76.76 0 0 1-.75-.75V5a.75.75 0 1 1 1.5 0v14a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgPause;
