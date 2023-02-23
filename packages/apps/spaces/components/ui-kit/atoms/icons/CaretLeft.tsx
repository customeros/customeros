import * as React from 'react';
import { SVGProps } from 'react';
const SvgCaretLeft = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M16 18.75a.74.74 0 0 1-.45-.15l-8-6a.75.75 0 0 1 0-1.2l8-6a.75.75 0 0 1 .79-.07.76.76 0 0 1 .41.67v12a.76.76 0 0 1-.41.67.84.84 0 0 1-.34.08ZM9.25 12l6 4.5v-9l-6 4.5Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCaretLeft;
