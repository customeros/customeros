import * as React from 'react';
import { SVGProps } from 'react';
const SvgSupplier = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M3.5 21.5a1.89 1.89 0 1 0 0-3.78 1.89 1.89 0 0 0 0 3.78ZM12.97 21.5a1.89 1.89 0 1 0 0-3.78 1.89 1.89 0 0 0 0 3.78Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M14.87 11.08v8.53a1.902 1.902 0 0 0-2.654-1.88 1.901 1.901 0 0 0-1.136 1.88H5.39a1.89 1.89 0 0 0-3.78 0v-5.69h8.52A2.849 2.849 0 0 1 13 11.08h1.87Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M14.87 11.08H13a2.85 2.85 0 0 0-2.84 2.84H1.61V3.5h5.68a7.58 7.58 0 0 1 7.58 7.58ZM4.45 10.13v3.79M11.5 11.5 9.66 9.66M11.08 8.24l-2.37 2.37M18.66 13.92h-3.79v5.68h3.79v-5.68ZM18.66 19.61V3.5M23.39 19.61h-4.73'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgSupplier;
