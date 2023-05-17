import * as React from 'react';
import { SVGProps } from 'react';
const SvgTicket = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <path
      d='M19 18.75H5A1.76 1.76 0 0 1 3.25 17v-2.5a.76.76 0 0 1 .75-.75A1.75 1.75 0 0 0 5.75 12 1.75 1.75 0 0 0 4 10.25a.76.76 0 0 1-.75-.75V7A1.76 1.76 0 0 1 5 5.25h14A1.76 1.76 0 0 1 20.75 7v2.5a.76.76 0 0 1-.75.75A1.75 1.75 0 0 0 18.25 12 1.75 1.75 0 0 0 20 13.75a.76.76 0 0 1 .75.75V17A1.76 1.76 0 0 1 19 18.75zM4.75 15.16V17a.25.25 0 0 0 .25.25h14a.25.25 0 0 0 .25-.25v-1.84A3.25 3.25 0 0 1 16.76 12a3.25 3.25 0 0 1 2.491-3.16V7a.25.25 0 0 0-.25-.25H5a.25.25 0 0 0-.25.25v1.84A3.25 3.25 0 0 1 7.24 12a3.25 3.25 0 0 1-2.49 3.16z'
      fill='currentColor'
    />
  </svg>
);
export default SvgTicket;
