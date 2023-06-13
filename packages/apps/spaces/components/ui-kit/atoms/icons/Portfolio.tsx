import * as React from 'react';
import { SVGProps } from 'react';
const SvgPortfolio = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M20.59 6.27H3.41A1.91 1.91 0 0 0 1.5 8.18v11.45a1.91 1.91 0 0 0 1.91 1.91h17.18a1.91 1.91 0 0 0 1.91-1.91V8.18a1.91 1.91 0 0 0-1.91-1.91Z'
      stroke='#343A40'
      strokeMiterlimit={10}
      strokeLinecap='square'
    />
    <path
      d='M13.91 13h4.77a3.812 3.812 0 0 0 3.82-3.86v-1a1.91 1.91 0 0 0-1.91-1.91H3.41A1.91 1.91 0 0 0 1.5 8.18v1A3.81 3.81 0 0 0 5.32 13h8.59Z'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
    <path
      d='M12 12v1.91M15.82 6.27H8.18l.96-3.81h5.72l.96 3.81Z'
      stroke='#343A40'
      strokeMiterlimit={10}
      strokeLinecap='square'
    />
  </svg>
);
export default SvgPortfolio;
