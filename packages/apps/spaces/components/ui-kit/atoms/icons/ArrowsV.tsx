import * as React from 'react';
import { SVGProps } from 'react';
const SvgArrowsV = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12.29 20.69a.78.78 0 0 0 .24-.16l4-4a.75.75 0 0 0-1.06-1.06l-2.72 2.72V5.81l2.72 2.72a.75.75 0 0 0 1.06-1.06l-4-4a.779.779 0 0 0-.53-.22.781.781 0 0 0-.53.22l-4 4a.75.75 0 0 0 1.06 1.06l2.72-2.72v12.38l-2.72-2.72a.75.75 0 0 0-1.06 1.06l4 4a.782.782 0 0 0 .53.22c.1 0 .198-.02.29-.06z'
      fill='currentColor'
    />
  </svg>
);
export default SvgArrowsV;
