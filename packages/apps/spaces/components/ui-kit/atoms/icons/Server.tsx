import * as React from 'react';
import { SVGProps } from 'react';
const SvgServer = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    xmlnsXlink='http://www.w3.org/1999/xlink'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M20.5 8.5v-3a1 1 0 0 0-1-1h-15a1 1 0 0 0-1 1v3a1 1 0 0 0 1 1 1 1 0 0 0-1 1v3a1 1 0 0 0 1 1 1 1 0 0 0-1 1v3a1 1 0 0 0 1 1h15a1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1 1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1 1 1 0 0 0 1-1zm-1 10h-15v-3h15v3zm0-5h-15v-3h15v3zm0-5h-15v-3h15v3z' />
      <use xlinkHref='#server_svg__a' />
      <use xlinkHref='#server_svg__a' x={2.5} />
      <use xlinkHref='#server_svg__a' y={5} />
      <use xlinkHref='#server_svg__a' x={2.5} y={5} />
      <use xlinkHref='#server_svg__a' y={10} />
      <use xlinkHref='#server_svg__a' x={2.5} y={10} />
    </g>
    <defs>
      <path
        id='server_svg__a'
        d='M6.25 7.75a.75.75 0 1 0 0-1.5.75.75 0 1 0 0 1.5z'
      />
    </defs>
  </svg>
);
export default SvgServer;
