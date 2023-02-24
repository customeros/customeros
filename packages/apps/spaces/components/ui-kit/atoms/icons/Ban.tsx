import * as React from 'react';
import { SVGProps } from 'react';
const SvgBan = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Zm-7.5 9a7.44 7.44 0 0 1 1.7-4.74L16.74 17.8A7.491 7.491 0 0 1 4.5 12Zm13.3 4.74L7.26 6.2A7.49 7.49 0 0 1 17.8 16.74Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgBan;
