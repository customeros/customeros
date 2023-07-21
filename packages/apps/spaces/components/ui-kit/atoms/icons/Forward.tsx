import * as React from 'react';
import { SVGProps } from 'react';
const SvgForward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 14 14'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M1.667 12.334V11.4c0-2.24 0-3.36.436-4.216A4 4 0 0 1 3.85 5.436C4.706 5 5.826 5 8.067 5h4.266m0 0L9 8.334M12.333 5 9 1.667'
      stroke='currentColor'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgForward;
