import * as React from 'react';
import { SVGProps } from 'react';
const SvgCaretRight = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M8 18.75a.76.76 0 0 1-.75-.75V6a.75.75 0 0 1 .42-.67.74.74 0 0 1 .78.07l8 6a.75.75 0 0 1 0 1.2l-8 6a.74.74 0 0 1-.45.15ZM8.75 7.5v9l6-4.5-6-4.5Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCaretRight;
