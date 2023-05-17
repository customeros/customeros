import * as React from 'react';
import { SVGProps } from 'react';
const SvgForward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M4 19.75a.8.8 0 0 1-.3-.06.76.76 0 0 1-.45-.69V5a.75.75 0 0 1 1.26-.55l7.47 7a.75.75 0 0 1 0 1.1l-7.47 7a.74.74 0 0 1-.51.2zm.75-13v10.52L10.37 12 4.75 6.75z' />
    <path d='M12.53 19.75a.72.72 0 0 1-.29-.06.75.75 0 0 1-.46-.69V5a.76.76 0 0 1 1.27-.55l7.46 7a.75.75 0 0 1 0 1.1l-7.46 7a.79.79 0 0 1-.52.2zm.75-13v10.52L18.9 12l-5.62-5.25z' />
  </svg>
);
export default SvgForward;
