import * as React from 'react';
import { SVGProps } from 'react';
const SvgReply = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={14}
    height={14}
    fill='currentColor'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M9.304 10.242a.677.677 0 0 1-.497-1.156l3.54-3.513-3.54-3.514a.677.677 0 0 1 .957-.957l4.029 3.992a.678.678 0 0 1 0 .957l-4.029 3.993a.633.633 0 0 1-.46.198Z'
      fill='#666'
    />
    <path
      d='M.677 13.096A.686.686 0 0 1 0 12.42V5.573a.686.686 0 0 1 .677-.678h12.646a.678.678 0 0 1 0 1.355H1.355v6.169a.686.686 0 0 1-.678.677Z'
      fill='#666'
    />
  </svg>
);
export default SvgReply;
