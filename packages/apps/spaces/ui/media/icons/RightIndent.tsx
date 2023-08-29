import * as React from 'react';
import { SVGProps } from 'react';
const SvgRightIndent = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M14 2.667H2m12 10.666H2m6-7.166H2m6 3.666H2M13.147 5.64l-2.578 1.933c-.193.145-.29.217-.324.306a.334.334 0 0 0 0 .242c.034.089.13.16.324.306l2.578 1.933c.274.206.412.309.527.307a.333.333 0 0 0 .255-.128c.071-.09.071-.262.071-.606V6.067c0-.344 0-.515-.07-.606a.333.333 0 0 0-.256-.128c-.115-.002-.253.101-.527.307Z'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgRightIndent;
