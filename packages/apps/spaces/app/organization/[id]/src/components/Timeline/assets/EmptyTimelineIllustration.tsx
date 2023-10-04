import * as React from 'react';
import { SVGProps } from 'react';
const SvgEmptyTimelineIlustration = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={172}
    height={152}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <circle cx={86} cy={72} r={64} fill='#E9D7FE' />
    <circle cx={20} cy={28} r={6} fill='#F4EBFF' />
    <circle cx={17} cy={120} r={8} fill='#F4EBFF' />
    <circle cx={160} cy={44} r={8} fill='#F4EBFF' />
    <circle cx={149} cy={17} r={5} fill='#F4EBFF' />
    <g filter='url(#empty-timeline-ilustration_svg__a)'>
      <path
        d='m86 44.65 5.616-2.353-5.615 2.392v66.975l51.747-21.685V22.975l.252-.106-.252-.112v-.112l-.13.055L86.408 0 34 22.325l.248.106L34 89.163l52 22.461V44.649Z'
        fill='#F9FAFB'
      />
      <path
        d='M86 44.649v66.974l-52-22.46.248-66.732L86 44.649Z'
        fill='url(#empty-timeline-ilustration_svg__b)'
      />
      <path
        d='M86.001 44.689v66.974l51.747-21.684V22.644L86.001 44.69Z'
        fill='url(#empty-timeline-ilustration_svg__c)'
      />
      <path
        d='m86 44.65 52-21.78L86.408 0 34 22.325l52 22.324Z'
        fill='#F4EBFF'
      />
      <path
        d='m53.602 13.975 51.802 22.548.47 16.43 14.494-5.98-.438-16.534L66.595 8.44l-12.993 5.535Z'
        fill='#E9D7FE'
      />
    </g>
    <defs>
      <linearGradient
        id='empty-timeline-ilustration_svg__b'
        x1={34}
        y1={89.499}
        x2={54.536}
        y2={33.882}
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#E9D7FE' />
        <stop offset={1} stopColor='#F1E6FF' />
      </linearGradient>
      <linearGradient
        id='empty-timeline-ilustration_svg__c'
        x1={86}
        y1={46}
        x2={144.635}
        y2={57.673}
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#E9D7FE' />
        <stop offset={1} stopColor='#F4EBFF' />
      </linearGradient>
      <filter
        id='empty-timeline-ilustration_svg__a'
        x={14}
        y={0}
        width={144}
        height={151.663}
        filterUnits='userSpaceOnUse'
        colorInterpolationFilters='sRGB'
      >
        <feFlood floodOpacity={0} result='BackgroundImageFix' />
        <feColorMatrix
          in='SourceAlpha'
          values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
          result='hardAlpha'
        />
        <feMorphology
          radius={4}
          in='SourceAlpha'
          result='effect1_dropShadow_541_521'
        />
        <feOffset dy={8} />
        <feGaussianBlur stdDeviation={4} />
        <feColorMatrix values='0 0 0 0 0.0627451 0 0 0 0 0.0941176 0 0 0 0 0.156863 0 0 0 0.03 0' />
        <feBlend in2='BackgroundImageFix' result='effect1_dropShadow_541_521' />
        <feColorMatrix
          in='SourceAlpha'
          values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
          result='hardAlpha'
        />
        <feMorphology
          radius={4}
          in='SourceAlpha'
          result='effect2_dropShadow_541_521'
        />
        <feOffset dy={20} />
        <feGaussianBlur stdDeviation={12} />
        <feColorMatrix values='0 0 0 0 0.0627451 0 0 0 0 0.0941176 0 0 0 0 0.156863 0 0 0 0.08 0' />
        <feBlend
          in2='effect1_dropShadow_541_521'
          result='effect2_dropShadow_541_521'
        />
        <feBlend
          in='SourceGraphic'
          in2='effect2_dropShadow_541_521'
          result='shape'
        />
      </filter>
    </defs>
  </svg>
);
export default SvgEmptyTimelineIlustration;
