import * as React from 'react';
import { SVGProps } from 'react';
const SvgCaretUp = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M18 16.75H6a.76.76 0 0 1-.67-.41.75.75 0 0 1 .07-.79l6-8a.77.77 0 0 1 1.2 0l6 8a.75.75 0 0 1-.6 1.2zm-10.5-1.5h9l-4.5-6-4.5 6z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCaretUp;
