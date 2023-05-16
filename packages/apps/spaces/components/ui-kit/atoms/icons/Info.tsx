import * as React from 'react';
import { SVGProps } from 'react';
const SvgInfo = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12 17.75a.76.76 0 0 1-.75-.75v-7a.75.75 0 1 1 1.5 0v7a.76.76 0 0 1-.75.75zm0-9.5a.76.76 0 0 1-.75-.75V7a.75.75 0 1 1 1.5 0v.5a.76.76 0 0 1-.75.75z' />
  </svg>
);
export default SvgInfo;
