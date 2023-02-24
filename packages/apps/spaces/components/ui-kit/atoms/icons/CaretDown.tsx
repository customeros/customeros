import * as React from 'react';
import { SVGProps } from 'react';
const SvgCaretDown = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 16.75a.74.74 0 0 1-.6-.3l-6-8a.75.75 0 0 1 .6-1.2h12a.76.76 0 0 1 .67.41.75.75 0 0 1-.07.79l-6 8a.74.74 0 0 1-.6.3Zm-4.5-8 4.5 6 4.5-6h-9Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCaretDown;
