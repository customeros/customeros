import * as React from 'react';
import { SVGProps } from 'react';
const SvgTwitter = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19.83 8v.52A11.41 11.41 0 0 1 8.35 20a11.41 11.41 0 0 1-6.2-1.81h1a8.09 8.09 0 0 0 5-1.73 4 4 0 0 1-3.78-2.8 4.662 4.662 0 0 0 1.83-.07A4 4 0 0 1 3 9.67a4.13 4.13 0 0 0 1.82.51 4.06 4.06 0 0 1-1.28-5.41A11.47 11.47 0 0 0 11.85 9a4.718 4.718 0 0 1-.1-.92 4 4 0 0 1 7-2.77 7.929 7.929 0 0 0 2.56-1 4 4 0 0 1-1.78 2.22 7.94 7.94 0 0 0 2.33-.62 8.908 8.908 0 0 1-2 2.09h-.03Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgTwitter;
