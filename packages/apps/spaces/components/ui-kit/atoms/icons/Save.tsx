import * as React from 'react';
import { SVGProps } from 'react';
const SvgSave = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M17 20.75H7A2.75 2.75 0 0 1 4.25 18V6A2.75 2.75 0 0 1 7 3.25h7.5c.199 0 .39.08.53.22L19.53 8c.14.14.22.331.22.53V18A2.75 2.75 0 0 1 17 20.75Zm-10-16A1.25 1.25 0 0 0 5.75 6v12A1.25 1.25 0 0 0 7 19.25h10A1.25 1.25 0 0 0 18.25 18V8.81l-4.06-4.06H7Z' />
      <path d='M16.75 20h-1.5v-6.25h-6.5V20h-1.5v-6.5a1.25 1.25 0 0 1 1.25-1.25h7a1.25 1.25 0 0 1 1.25 1.25V20ZM12.47 8.75H8.53a1.29 1.29 0 0 1-1.28-1.3V4h1.5v3.25h3.5V4h1.5v3.45a1.29 1.29 0 0 1-1.28 1.3Z' />
    </g>
  </svg>
);
export default SvgSave;
