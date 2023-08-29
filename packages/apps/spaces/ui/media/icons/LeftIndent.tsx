import * as React from 'react';
import { SVGProps } from 'react';
const SvgLeftIndent = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M14 6.167H8m6-3.5H2m12 7.166H8m6 3.5H2m.853-7.626L5.431 7.64c.193.145.29.217.324.306.03.078.03.164 0 .242-.034.088-.13.16-.324.305l-2.578 1.934c-.274.206-.412.309-.527.306a.333.333 0 0 1-.255-.128C2 10.515 2 10.343 2 10V6.133c0-.343 0-.515.07-.605a.333.333 0 0 1 .256-.128c.115-.002.253.1.527.307Z'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgLeftIndent;
