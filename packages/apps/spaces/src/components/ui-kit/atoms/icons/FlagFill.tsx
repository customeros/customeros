import * as React from 'react';
import { SVGProps } from 'react';
const SvgFlagFill = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19.42 4.45a.74.74 0 0 0-.69-.08 13.18 13.18 0 0 1-3.73 1 8.46 8.46 0 0 1-2.3-1 8.76 8.76 0 0 0-3-1.16c-1.29-.12-4.36.89-5 1.1a.75.75 0 0 0-.51.71V20a.75.75 0 1 0 1.5 0v-5.86a15.998 15.998 0 0 1 3.86-.85 8.469 8.469 0 0 1 2.4 1 9.11 9.11 0 0 0 2.82 1.13H15a16.369 16.369 0 0 0 4.21-1.13.76.76 0 0 0 .48-.7V5.07a.74.74 0 0 0-.27-.62Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgFlagFill;
