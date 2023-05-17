import * as React from 'react';
import { SVGProps } from 'react';
const SvgRefresh = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12 20.75a7.25 7.25 0 0 1-7.25-7.25A7.25 7.25 0 0 1 12 6.25h2.5a.75.75 0 0 1 .75.75.75.75 0 0 1-.75.75H12a5.75 5.75 0 0 0-5.312 3.55 5.75 5.75 0 0 0 1.246 6.266 5.75 5.75 0 0 0 6.266 1.246 5.75 5.75 0 0 0 3.55-5.312.75.75 0 0 1 .75-.75.75.75 0 0 1 .75.75A7.26 7.26 0 0 1 12 20.75z' />
    <path d='M12 10.75a.74.74 0 0 1-.53-.22.75.75 0 0 1 0-1.06L13.94 7l-2.47-2.47a.75.75 0 0 1-.239-.535.75.75 0 0 1 .22-.544.75.75 0 0 1 .544-.22.75.75 0 0 1 .535.239l3 3a.75.75 0 0 1 0 1.06l-3 3a.74.74 0 0 1-.53.22z' />
  </svg>
);
export default SvgRefresh;
