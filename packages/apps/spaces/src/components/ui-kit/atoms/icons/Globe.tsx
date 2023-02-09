import * as React from 'react';
import { SVGProps } from 'react';
const SvgGlobe = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Zm7.46 8.25H16.7a13 13 0 0 0-2.94-6.53 7.52 7.52 0 0 1 5.7 6.53Zm-10.65 1.5h6.38A13.18 13.18 0 0 1 12 19.1a13.18 13.18 0 0 1-3.19-6.35Zm0-1.5A13.18 13.18 0 0 1 12 4.9a13.18 13.18 0 0 1 3.19 6.35H8.81Zm1.43-6.53a13 13 0 0 0-2.94 6.53H4.54a7.52 7.52 0 0 1 5.7-6.53Zm-5.7 8H7.3a13 13 0 0 0 2.94 6.53 7.52 7.52 0 0 1-5.7-6.5v-.03Zm9.22 6.53a13 13 0 0 0 2.94-6.53h2.76a7.52 7.52 0 0 1-5.7 6.56v-.03Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgGlobe;
