import * as React from 'react';
import { SVGProps } from 'react';
const SvgUser = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12 12.25a3.75 3.75 0 0 1-2.652-6.402A3.75 3.75 0 0 1 15.75 8.5 3.75 3.75 0 0 1 12 12.25zm0-6a2.25 2.25 0 0 0-1.591 3.841A2.25 2.25 0 0 0 14.25 8.5 2.25 2.25 0 0 0 12 6.25zm7 13a.76.76 0 0 1-.75-.75c0-1.95-1.06-3.25-6.25-3.25s-6.25 1.3-6.25 3.25a.75.75 0 1 1-1.5 0c0-4.75 5.43-4.75 7.75-4.75s7.75 0 7.75 4.75a.76.76 0 0 1-.75.75z' />
  </svg>
);
export default SvgUser;
