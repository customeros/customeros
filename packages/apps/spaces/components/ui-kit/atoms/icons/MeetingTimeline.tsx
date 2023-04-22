import * as React from 'react';
import { SVGProps } from 'react';
const SvgMeetingTimeline = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 142 60'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g filter='url(#meeting-timeline_svg__a)'>
      <path
        d='M5 9v34a8 8 0 0 0 8 8h116a8 8 0 0 0 8-8V9a8 8 0 0 0-8-8H93.809a8 8 0 0 0-7.372 4.892l-.09.215A8 8 0 0 1 78.974 11H62.14a8 8 0 0 1-7.312-4.754l-.218-.49A8 8 0 0 0 47.298 1H13a8 8 0 0 0-8 8Z'
        fill='#FBFBFB'
      />
      <path
        d='M5 9v34a8 8 0 0 0 8 8h116a8 8 0 0 0 8-8V9a8 8 0 0 0-8-8H93.809a8 8 0 0 0-7.372 4.892l-.09.215A8 8 0 0 1 78.974 11H62.14a8 8 0 0 1-7.312-4.754l-.218-.49A8 8 0 0 0 47.298 1H13a8 8 0 0 0-8 8Z'
        stroke='#3987A6'
        strokeWidth={0.5}
      />
    </g>
    <defs>
      <filter
        id='meeting-timeline_svg__a'
        x={0.75}
        y={0.75}
        width={140.5}
        height={58.5}
        filterUnits='userSpaceOnUse'
        colorInterpolationFilters='sRGB'
      >
        <feFlood floodOpacity={0} result='BackgroundImageFix' />
        <feColorMatrix
          in='SourceAlpha'
          values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
          result='hardAlpha'
        />
        <feOffset dy={4} />
        <feGaussianBlur stdDeviation={2} />
        <feComposite in2='hardAlpha' operator='out' />
        <feColorMatrix values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.25 0' />
        <feBlend
          in2='BackgroundImageFix'
          result='effect1_dropShadow_737_8843'
        />
        <feBlend
          in='SourceGraphic'
          in2='effect1_dropShadow_737_8843'
          result='shape'
        />
      </filter>
    </defs>
  </svg>
);
export default SvgMeetingTimeline;
