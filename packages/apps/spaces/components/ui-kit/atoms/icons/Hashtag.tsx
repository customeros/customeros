import * as React from 'react';
import { SVGProps } from 'react';
const SvgHashtag = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M21 7.25h-3l.77-3.07a.752.752 0 0 0-1.46-.36l-.86 3.43H10l.77-3.07a.752.752 0 0 0-1.46-.36l-.9 3.43H5a.75.75 0 1 0 0 1.5h3l-1.63 6.5H3a.75.75 0 1 0 0 1.5h3l-.77 3.07a.752.752 0 1 0 1.46.36l.86-3.43H14l-.77 3.07a.752.752 0 0 0 1.46.36l.86-3.43H19a.75.75 0 1 0 0-1.5h-3l1.63-6.5H21a.75.75 0 1 0 0-1.5Zm-5 1.5-1.63 6.5H8l1.63-6.5H16Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgHashtag;
