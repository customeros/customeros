import * as React from 'react';
import { SVGProps } from 'react';
const SvgSlack = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M7.36 14.11a1.68 1.68 0 1 1-1.68-1.68h1.68v1.68Zm.85 0a1.68 1.68 0 1 1 3.36 0v4.21a1.68 1.68 0 1 1-3.36 0v-4.21Zm1.68-6.75a1.68 1.68 0 1 1 1.68-1.68v1.68H9.89Zm0 .85a1.68 1.68 0 1 1 0 3.36H5.68a1.68 1.68 0 1 1 0-3.36h4.21Zm6.75 1.68a1.68 1.68 0 1 1 1.68 1.68h-1.68V9.89Zm-.85 0a1.68 1.68 0 1 1-3.36 0V5.68a1.68 1.68 0 1 1 3.36 0v4.21Zm-1.68 6.75a1.68 1.68 0 1 1-1.68 1.68v-1.68h1.68Zm0-.85a1.68 1.68 0 1 1 0-3.36h4.21a1.68 1.68 0 1 1 0 3.36h-4.21Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgSlack;
