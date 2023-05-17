import * as React from 'react';
import { SVGProps } from 'react';
const SvgPlus = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <path
      d='M12.75 11.25V5a.75.75 0 0 0-.75-.75.75.75 0 0 0-.75.75v6.25H5a.75.75 0 0 0-.75.75.75.75 0 0 0 .75.75h6.25V19a.76.76 0 0 0 .75.75.75.75 0 0 0 .75-.75v-6.25H19a.75.75 0 0 0 .75-.75.76.76 0 0 0-.75-.75h-6.25z'
      fill='currentColor'
    />
  </svg>
);
export default SvgPlus;
