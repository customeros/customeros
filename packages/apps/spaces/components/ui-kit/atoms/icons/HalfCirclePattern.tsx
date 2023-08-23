import * as React from 'react';
import { SVGProps } from 'react';
const SvgHalfCirclePattern = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 424 339'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <mask
      id='half-circle-pattern_svg__b'
      style={{
        maskType: 'alpha',
      }}
      maskUnits='userSpaceOnUse'
      x={0}
      y={0}
      width={424}
      height={339}
    >
      <path fill='url(#half-circle-pattern_svg__a)' d='M0 0h424v339H0z' />
    </mask>
    <g mask='url(#half-circle-pattern_svg__b)' stroke='#EAECF0'>
      <circle cx={212} cy={99} r={47.5} />
      <circle cx={212} cy={99} r={79.5} />
      <circle cx={212} cy={99} r={111.5} />
      <circle cx={212} cy={99} r={143.5} />
      <circle cx={212} cy={99} r={143.5} />
      <circle cx={212} cy={99} r={175.5} />
      <circle cx={212} cy={99} r={207.5} />
      <circle cx={212} cy={99} r={239.5} />
    </g>
    <defs>
      <radialGradient
        id='half-circle-pattern_svg__a'
        cx={0}
        cy={0}
        r={1}
        gradientUnits='userSpaceOnUse'
        gradientTransform='matrix(0 -339 221.262 0 212 339)'
      >
        <stop />
        <stop offset={0.958} stopOpacity={0} />
      </radialGradient>
    </defs>
  </svg>
);
export default SvgHalfCirclePattern;
