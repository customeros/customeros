import * as React from 'react';
import { SVGProps } from 'react';
const SvgSync = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M18.43 4.25a.76.76 0 0 0-.75.75v2.43l-.84-.84a7.24 7.24 0 0 0-12 2.78.74.74 0 0 0 .46 1 .729.729 0 0 0 .25 0 .76.76 0 0 0 .71-.51 5.63 5.63 0 0 1 1.37-2.2 5.76 5.76 0 0 1 8.13 0l.84.84h-2.41a.75.75 0 1 0 0 1.5h4.24a.74.74 0 0 0 .75-.75V5a.75.75 0 0 0-.75-.75ZM18.68 13.68a.761.761 0 0 0-1 .47 5.63 5.63 0 0 1-1.37 2.2 5.76 5.76 0 0 1-8.13 0l-.84-.84h2.47a.75.75 0 1 0 0-1.5H5.57a.74.74 0 0 0-.75.75V19a.75.75 0 1 0 1.5 0v-2.43l.84.84a7.24 7.24 0 0 0 12-2.78.74.74 0 0 0-.48-.95Z' />
    </g>
  </svg>
);
export default SvgSync;
