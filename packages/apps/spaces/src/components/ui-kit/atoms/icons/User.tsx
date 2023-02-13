import * as React from 'react';
import { SVGProps } from 'react';
const SvgUser = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 12.25a3.75 3.75 0 1 1 0-7.5 3.75 3.75 0 0 1 0 7.5Zm0-6a2.25 2.25 0 1 0 0 4.5 2.25 2.25 0 0 0 0-4.5ZM19 19.25a.76.76 0 0 1-.75-.75c0-1.95-1.06-3.25-6.25-3.25-5.19 0-6.25 1.3-6.25 3.25a.75.75 0 1 1-1.5 0c0-4.75 5.43-4.75 7.75-4.75 2.32 0 7.75 0 7.75 4.75a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgUser;
