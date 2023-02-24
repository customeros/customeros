import * as React from 'react';
import { SVGProps } from 'react';
const SvgQrcode = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12.38 3.88v7h7v-7h-7Zm5.5 5.5h-4v-4h4v4ZM3.88 10.88h7v-7h-7v7Zm1.5-5.5h4v4h-4v-4ZM3.88 19.38h7v-7h-7v7Zm1.5-5.5h4v4h-4v-4ZM14.13 12.38h-1.75v1.75h1.75v-1.75ZM17.63 12.38h-1.75v1.75h1.75v-1.75ZM15.88 14.13h-1.75v1.75h1.75v-1.75ZM19.38 14.13h-1.75v1.75h1.75v-1.75ZM14.13 15.88h-1.75v1.75h1.75v-1.75ZM17.63 15.88h-1.75v1.75h1.75v-1.75Z' />
      <path d='M15.88 17.63h-1.75v1.75h1.75v-1.75ZM19.38 17.63h-1.75v1.75h1.75v-1.75Z' />
    </g>
  </svg>
);
export default SvgQrcode;
