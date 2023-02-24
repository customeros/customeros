import * as React from 'react';
import { SVGProps } from 'react';
const SvgIdCard = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M19 18.75H5A1.76 1.76 0 0 1 3.25 17V7A1.76 1.76 0 0 1 5 5.25h14A1.76 1.76 0 0 1 20.75 7v10A1.76 1.76 0 0 1 19 18.75Zm-14-12a.25.25 0 0 0-.25.25v10a.25.25 0 0 0 .25.25h14a.25.25 0 0 0 .25-.25V7a.25.25 0 0 0-.25-.25H5Z' />
      <path d='M9 11.75a2 2 0 1 1 0-4 2 2 0 0 1 0 4Zm0-2.5a.5.5 0 1 0 0 1 .5.5 0 0 0 0-1ZM12 15.75a.76.76 0 0 1-.75-.75c0-.68-.17-1.25-2.25-1.25s-2.25.57-2.25 1.25a.75.75 0 1 1-1.5 0c0-2.75 2.82-2.75 3.75-2.75.93 0 3.75 0 3.75 2.75a.76.76 0 0 1-.75.75ZM17 10.75h-3a.75.75 0 1 1 0-1.5h3a.75.75 0 1 1 0 1.5ZM16 13.75h-2a.75.75 0 1 1 0-1.5h2a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgIdCard;
