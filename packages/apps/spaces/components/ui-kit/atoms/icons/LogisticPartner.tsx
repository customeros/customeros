import * as React from 'react';
import { SVGProps } from 'react';
const SvgLogisticPartner = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M6.27 19.64a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM17.73 19.64a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M4.36 17.73H1.5V4.36h17.18v3.82L20.59 12l1.91 1.07v4.66h-2.86M15.82 17.73H8.18'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path d='M20.59 12h-5.73V8.18h3.82' stroke='#000' strokeMiterlimit={10} />
  </svg>
);
export default SvgLogisticPartner;
