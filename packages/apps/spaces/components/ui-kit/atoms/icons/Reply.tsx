import * as React from 'react';
import { SVGProps } from 'react';
const SvgReply = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 14 14'
    fill='#666'
    {...props}
  >
    <path d='M9.304 10.242a.683.683 0 0 1-.64-.41.685.685 0 0 1-.043-.395.68.68 0 0 1 .186-.351l3.541-3.514-3.541-3.514a.677.677 0 0 1 .958-.958l4.028 3.992a.68.68 0 0 1 0 .957l-4.028 3.992a.628.628 0 0 1-.461.199z' />
    <path d='M.677 13.096A.689.689 0 0 1 0 12.419V5.573a.687.687 0 0 1 .677-.677h12.645a.68.68 0 0 1 .677.677.68.68 0 0 1-.677.677H1.355v6.169a.687.687 0 0 1-.677.677z' />
  </svg>
);
export default SvgReply;
