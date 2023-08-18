import * as React from 'react';
import { SVGProps } from 'react';
const SvgExternalLink = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 14 14'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M13 5V1m0 0H9m4 0L7 7M5.667 1H4.2c-1.12 0-1.68 0-2.108.218a2 2 0 0 0-.874.874C1 2.52 1 3.08 1 4.2v5.6c0 1.12 0 1.68.218 2.108a2 2 0 0 0 .874.874C2.52 13 3.08 13 4.2 13h5.6c1.12 0 1.68 0 2.108-.218a2 2 0 0 0 .874-.874C13 11.48 13 10.92 13 9.8V8.333'
      stroke='currentColor'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgExternalLink;
