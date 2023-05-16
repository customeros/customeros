import * as React from 'react';
import { SVGProps } from 'react';
const SvgQrcode = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12.38 3.88v7h7v-7h-7zm5.5 5.5h-4v-4h4v4zm-14 1.5h7v-7h-7v7zm1.5-5.5h4v4h-4v-4zm-1.5 14h7v-7h-7v7zm1.5-5.5h4v4h-4v-4zm8.75-1.5h-1.75v1.75h1.75v-1.75zm3.5 0h-1.75v1.75h1.75v-1.75zM15.88 14.13h-1.75v1.75h1.75v-1.75zm3.5 0h-1.75v1.75h1.75v-1.75zM14.13 15.88h-1.75v1.75h1.75v-1.75zm3.5 0h-1.75v1.75h1.75v-1.75zM15.88 17.63h-1.75v1.75h1.75v-1.75zm3.5 0h-1.75v1.75h1.75v-1.75z' />
  </svg>
);
export default SvgQrcode;
