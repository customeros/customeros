import * as React from 'react';
import { SVGProps } from 'react';
const SvgBolt = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M11.09 21.5a.672.672 0 0 1-.24 0 .83.83 0 0 1-.59-.81v-.11l.9-6.35H6.82a.8.8 0 0 1-.71-.43.85.85 0 0 1 0-.86l2-3.49 4.1-6.52a.79.79 0 0 1 .92-.35.83.83 0 0 1 .59.81v.11l-.9 6.35h4.35a.8.8 0 0 1 .71.43.85.85 0 0 1 0 .86l-2 3.49-4.1 6.52a.79.79 0 0 1-.69.35Zm-3.16-8.81h4a.84.84 0 0 1 .83.85v.11l-.59 4.14 2.5-4 1.44-2.48h-4a.839.839 0 0 1-.83-.85v-.11l.59-4.14-2.5 4-1.44 2.48Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgBolt;
