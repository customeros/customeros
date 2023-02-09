import * as React from 'react';
import { SVGProps } from 'react';
const SvgTicket = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19 18.75H5A1.76 1.76 0 0 1 3.25 17v-2.5a.76.76 0 0 1 .75-.75 1.75 1.75 0 0 0 0-3.5.76.76 0 0 1-.75-.75V7A1.76 1.76 0 0 1 5 5.25h14A1.76 1.76 0 0 1 20.75 7v2.5a.76.76 0 0 1-.75.75 1.75 1.75 0 0 0 0 3.5.76.76 0 0 1 .75.75V17A1.76 1.76 0 0 1 19 18.75ZM4.75 15.16V17a.25.25 0 0 0 .25.25h14a.25.25 0 0 0 .25-.25v-1.84a3.25 3.25 0 0 1 0-6.32V7a.25.25 0 0 0-.25-.25H5a.25.25 0 0 0-.25.25v1.84a3.25 3.25 0 0 1 0 6.32Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgTicket;
