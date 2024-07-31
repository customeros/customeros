import { SVGProps } from 'react';
const SvgHalfCirclePattern = (props: SVGProps<SVGSVGElement>) => (
  <svg
    fill='none'
    viewBox='0 0 424 339'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <mask
      x={0}
      y={0}
      width={424}
      height={339}
      maskUnits='userSpaceOnUse'
      id='half-circle-pattern_svg__b'
      style={{
        maskType: 'alpha',
      }}
    >
      <path d='M0 0h424v339H0z' fill='url(#half-circle-pattern_svg__a)' />
    </mask>
    <g stroke='#EAECF0' mask='url(#half-circle-pattern_svg__b)'>
      <circle cy={99} cx={212} r={47.5} />
      <circle cy={99} cx={212} r={79.5} />
      <circle cy={99} cx={212} r={111.5} />
      <circle cy={99} cx={212} r={143.5} />
      <circle cy={99} cx={212} r={143.5} />
      <circle cy={99} cx={212} r={175.5} />
      <circle cy={99} cx={212} r={207.5} />
      <circle cy={99} cx={212} r={239.5} />
    </g>
    <defs>
      <radialGradient
        r={1}
        cx={0}
        cy={0}
        gradientUnits='userSpaceOnUse'
        id='half-circle-pattern_svg__a'
        gradientTransform='matrix(0 -339 221.262 0 212 339)'
      >
        <stop />
        <stop offset={0.958} stopOpacity={0} />
      </radialGradient>
    </defs>
  </svg>
);
export default SvgHalfCirclePattern;
