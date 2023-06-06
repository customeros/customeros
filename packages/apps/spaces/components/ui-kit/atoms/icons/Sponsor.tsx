import * as React from 'react';
import { SVGProps } from 'react';
const SvgSponsor = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M15.82 14.86H1.5V1.5h21v13.36h-2.86'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M9.14 11.05h3.34a1.44 1.44 0 0 0 1.011-2.451 1.43 1.43 0 0 0-1.01-.419h-1a1.43 1.43 0 0 1 0-2.86h3.34M12 3.41v1.91M12 11.05v1.9M17.73 22.5a6.52 6.52 0 0 0 1.91-4.61v-4A1.91 1.91 0 0 0 17.73 12a1.9 1.9 0 0 0-1.91 1.91v2.86l-.2.1a3.1 3.1 0 0 0-1.71 2.77M8.18 14.86v3a6.52 6.52 0 0 0 1.91 4.61'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgSponsor;
