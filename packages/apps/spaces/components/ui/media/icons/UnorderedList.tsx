import * as React from 'react';
import { SVGProps } from 'react';
const SvgUnorderedList = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M14 3.333H6.667M14 12.667H6.667M14 8H6.667M4 3.333a1 1 0 1 1-2 0 1 1 0 0 1 2 0Zm0 9.334a1 1 0 1 1-2 0 1 1 0 0 1 2 0ZM4 8a1 1 0 1 1-2 0 1 1 0 0 1 2 0Z'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgUnorderedList;
