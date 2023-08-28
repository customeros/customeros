import * as React from 'react';
import { SVGProps } from 'react';
const SvgHeading1 = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      fillRule='evenodd'
      clipRule='evenodd'
      d='M1.5 3.474c0-.262.212-.474.474-.474h2.013a.474.474 0 0 1 0 .947h-.533v3.08h5.092v-3.08h-.533a.474.474 0 0 1 0-.947h2.013a.474.474 0 0 1 0 .947h-.533v7.106h.533a.474.474 0 0 1 0 .947H8.013a.474.474 0 0 1 0-.947h.533v-3.08H3.454v3.08h.533a.474.474 0 0 1 0 .947H1.974a.474.474 0 0 1 0-.947h.533V3.947h-.533a.474.474 0 0 1-.474-.473Z'
      fill='currentColor'
    />
    <path
      d='M12 7.179c0-.2.12-.352.28-.437l.944-.512c.296-.161.456-.23.688-.23.344 0 .528.237.528.49v4.967a.558.558 0 0 1-.568.543.558.558 0 0 1-.568-.543V7.3l-.576.298a.571.571 0 0 1-.232.054A.488.488 0 0 1 12 7.18Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgHeading1;
