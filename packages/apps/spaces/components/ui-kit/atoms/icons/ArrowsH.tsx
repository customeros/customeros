import * as React from 'react';
import { SVGProps } from 'react';
const SvgArrowsH = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M20.69 11.71a.779.779 0 0 0-.16-.24l-4-4a.749.749 0 1 0-1.06 1.06l2.72 2.72H5.81l2.72-2.72a.75.75 0 0 0-1.06-1.06l-4 4a.78.78 0 0 0-.22.53.78.78 0 0 0 .22.53l4 4a.75.75 0 0 0 1.06-1.06l-2.72-2.72h12.38l-2.72 2.72a.75.75 0 0 0 1.06 1.06l4-4a.779.779 0 0 0 .22-.53.73.73 0 0 0-.06-.29Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgArrowsH;
