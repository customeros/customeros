import * as React from 'react';
import { SVGProps } from 'react';
const SvgLicensingPartner = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='m12 10.09 1.18 2.52 2.64.4-1.91 1.95.45 2.77L12 16.42l-2.36 1.31.45-2.77-1.91-1.95 2.64-.4L12 10.09Z'
      fill='#000'
    />
    <path
      d='M15.82 2.46h4.77V22.5H3.41V2.46h4.77'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M15.82 1.5v1.91a1.92 1.92 0 0 1-1.91 1.91h-3.82a1.92 1.92 0 0 1-1.91-1.91V1.5h7.64Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgLicensingPartner;
