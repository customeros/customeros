import * as React from 'react';
import { SVGProps } from 'react';
const SvgCopy = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M8 8V5.2c0-1.12 0-1.68.218-2.108a2 2 0 0 1 .874-.874C9.52 2 10.08 2 11.2 2h7.6c1.12 0 1.68 0 2.108.218a2 2 0 0 1 .874.874C22 3.52 22 4.08 22 5.2v7.6c0 1.12 0 1.68-.218 2.108a2 2 0 0 1-.874.874C20.48 16 19.92 16 18.8 16H16M5.2 22h7.6c1.12 0 1.68 0 2.108-.218a2 2 0 0 0 .874-.874C16 20.48 16 19.92 16 18.8v-7.6c0-1.12 0-1.68-.218-2.108a2 2 0 0 0-.874-.874C14.48 8 13.92 8 12.8 8H5.2c-1.12 0-1.68 0-2.108.218a2 2 0 0 0-.874.874C2 9.52 2 10.08 2 11.2v7.6c0 1.12 0 1.68.218 2.108a2 2 0 0 0 .874.874C3.52 22 4.08 22 5.2 22Z'
      stroke='currentColor'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgCopy;
