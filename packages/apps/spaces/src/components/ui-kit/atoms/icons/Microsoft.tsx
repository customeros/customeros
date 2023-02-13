import * as React from 'react';
import { SVGProps } from 'react';
const SvgMicrosoft = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M4 4h7.5v7.5H4V4Zm8.5 0H20v7.5h-7.5V4ZM4 12.5h7.5V20H4v-7.5Zm8.5 0H20V20h-7.5v-7.5Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgMicrosoft;
