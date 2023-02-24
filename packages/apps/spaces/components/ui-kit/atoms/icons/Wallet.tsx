import * as React from 'react';
import { SVGProps } from 'react';
const SvgWallet = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M19 7.25h-.25V5A1.76 1.76 0 0 0 17 3.25a.67.67 0 0 0-.24 0l-11.9 4h-.27l-.17.06h-.14l-.16.09-.12.17-.14.12-.11.1-.12.15a.39.39 0 0 0-.08.1 1.62 1.62 0 0 0-.1.18l-.06.11a1.87 1.87 0 0 0-.07.22.45.45 0 0 1 0 .11c-.01.113-.01.227 0 .34v10A1.76 1.76 0 0 0 5 20.75h14A1.76 1.76 0 0 0 20.75 19V9A1.76 1.76 0 0 0 19 7.25Zm-1.92-2.49a.26.26 0 0 1 .17.24v2.25H9.62l7.46-2.49ZM19.25 19a.25.25 0 0 1-.25.25H5a.25.25 0 0 1-.25-.25V9A.25.25 0 0 1 5 8.75h14a.25.25 0 0 1 .25.25v10Z' />
      <path d='M16.5 15.25a1.25 1.25 0 1 0 0-2.5 1.25 1.25 0 0 0 0 2.5Z' />
    </g>
  </svg>
);
export default SvgWallet;
