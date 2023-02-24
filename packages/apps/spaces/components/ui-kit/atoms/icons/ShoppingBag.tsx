import * as React from 'react';
import { SVGProps } from 'react';
const SvgShoppingBag = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19.5 8.25h-3v-.5a4.5 4.5 0 1 0-9 0v.5h-3A1.25 1.25 0 0 0 3.25 9.5V18A2.75 2.75 0 0 0 6 20.75h12A2.75 2.75 0 0 0 20.75 18V9.5a1.25 1.25 0 0 0-1.25-1.25ZM9 7.75a3 3 0 1 1 6 0v.5H9v-.5ZM19.25 18A1.25 1.25 0 0 1 18 19.25H6A1.25 1.25 0 0 1 4.75 18V9.75H7.5V12A.75.75 0 1 0 9 12V9.75h6V12a.75.75 0 1 0 1.5 0V9.75h2.75V18Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgShoppingBag;
