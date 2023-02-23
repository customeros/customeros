import * as React from 'react';
import { SVGProps } from 'react';
const SvgYoutube = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19.67 8.14a2 2 0 0 0-1.42-1.43A48.443 48.443 0 0 0 12 6.38a48.442 48.442 0 0 0-6.25.33 2 2 0 0 0-1.42 1.43A21.27 21.27 0 0 0 4 12c-.008 1.3.102 2.6.33 3.88a2 2 0 0 0 1.42 1.4c2.074.245 4.162.355 6.25.33a48.434 48.434 0 0 0 6.25-.33 2 2 0 0 0 1.42-1.4c.228-1.28.338-2.58.33-3.88a21.273 21.273 0 0 0-.33-3.86Zm-9.31 6.25V9.63L14.55 12l-4.19 2.38v.01Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgYoutube;
