import * as React from 'react';
import { SVGProps } from 'react';
const SvgForward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M4 19.75a.8.8 0 0 1-.3-.06.76.76 0 0 1-.45-.69V5a.75.75 0 0 1 1.26-.55l7.47 7a.752.752 0 0 1 0 1.1l-7.47 7a.74.74 0 0 1-.51.2Zm.75-13v10.52L10.37 12 4.75 6.75Z'
      fill='currentColor'
    />
    <path
      d='M12.53 19.75a.75.75 0 0 1-.75-.75V5a.76.76 0 0 1 1.27-.55l7.46 7a.752.752 0 0 1 0 1.1l-7.46 7a.79.79 0 0 1-.52.2Zm.75-13v10.52L18.9 12l-5.62-5.25Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgForward;
