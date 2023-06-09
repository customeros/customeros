import * as React from 'react';
import { SVGProps } from 'react';
const SvgOutsourcingProvider = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M6.27 22.5h7.57a4.77 4.77 0 0 0 3.38-1.4l4.71-4.71a2.025 2.025 0 1 0-2.86-2.87l-5.16 5.16H13a1.91 1.91 0 0 0 1.88-2.23 2 2 0 0 0-2-1.59H9.14l-.93-.46a4.66 4.66 0 0 0-6.273 2.163A4.57 4.57 0 0 0 1.5 18.57v.11M12.96 18.68H8.18M12 12.96A5.73 5.73 0 1 0 12 1.5a5.73 5.73 0 0 0 0 11.46Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M12 12.96c1.055 0 1.91-2.565 1.91-5.73 0-3.165-.855-5.73-1.91-5.73s-1.91 2.565-1.91 5.73c0 3.165.855 5.73 1.91 5.73ZM6.27 7.23h11.46'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgOutsourcingProvider;
