import { SVGProps } from 'react';

export const WelcomePage = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width='152'
    fill='none'
    height='131'
    viewBox='0 0 152 131'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <circle r='52' cx='76' cy='58' fill='#E9D7FE' />
    <circle r='5' cx='21' cy='25' fill='#F4EBFF' />
    <circle r='7' cx='18' cy='109' fill='#F4EBFF' />
    <circle r='7' cy='41' cx='145' fill='#F4EBFF' />
    <circle r='4' cy='14' cx='134' fill='#F4EBFF' />
    <g filter='url(#filter0_dd_632_3973)'>
      <path
        fill='#F9FAFB'
        d='M76 36.0629L80.536 34.1629L76.0006 36.0951V90.1898L117.796 72.6754V18.5566L118 18.4712L117.796 18.3809V18.29L117.692 18.3346L76.3298 0L34 18.0314L34.2004 18.1175L34 72.0166L76 90.1576L76 36.0629Z'
      />
      <path
        fill='url(#paint0_linear_632_3973)'
        d='M76 36.063V90.1577L34 72.0166L34.2004 18.1175L76 36.063Z'
      />
      <path
        fill='url(#paint1_linear_632_3973)'
        d='M76.001 36.0951V90.1898L117.797 72.6755V18.2901L76.001 36.0951Z'
      />
      <path
        fill='#F4EBFF'
        d='M76 36.0629L118 18.4712L76.3298 0L34 18.0314L76 36.0629Z'
      />
      <path
        fill='#E9D7FE'
        d='M49.832 11.2874L91.6718 29.4987L92.0519 42.7698L103.758 37.9395L103.404 24.5846L60.3266 6.81694L49.832 11.2874Z'
      />
    </g>
    <defs>
      <filter
        y='0'
        x='14'
        width='124'
        height='130.19'
        id='filter0_dd_632_3973'
        filterUnits='userSpaceOnUse'
        color-interpolation-filters='sRGB'
      >
        <feFlood flood-opacity='0' result='BackgroundImageFix' />
        <feColorMatrix
          type='matrix'
          in='SourceAlpha'
          result='hardAlpha'
          values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
        />
        <feMorphology
          radius='4'
          operator='erode'
          in='SourceAlpha'
          result='effect1_dropShadow_632_3973'
        />
        <feOffset dy='8' />
        <feGaussianBlur stdDeviation='4' />
        <feColorMatrix
          type='matrix'
          values='0 0 0 0 0.0627451 0 0 0 0 0.0941176 0 0 0 0 0.156863 0 0 0 0.03 0'
        />
        <feBlend
          mode='normal'
          in2='BackgroundImageFix'
          result='effect1_dropShadow_632_3973'
        />
        <feColorMatrix
          type='matrix'
          in='SourceAlpha'
          result='hardAlpha'
          values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
        />
        <feMorphology
          radius='4'
          operator='erode'
          in='SourceAlpha'
          result='effect2_dropShadow_632_3973'
        />
        <feOffset dy='20' />
        <feGaussianBlur stdDeviation='12' />
        <feColorMatrix
          type='matrix'
          values='0 0 0 0 0.0627451 0 0 0 0 0.0941176 0 0 0 0 0.156863 0 0 0 0.08 0'
        />
        <feBlend
          mode='normal'
          in2='effect1_dropShadow_632_3973'
          result='effect2_dropShadow_632_3973'
        />
        <feBlend
          mode='normal'
          result='shape'
          in='SourceGraphic'
          in2='effect2_dropShadow_632_3973'
        />
      </filter>
      <linearGradient
        x1='34'
        x2='50.587'
        y1='72.2885'
        y2='27.3665'
        id='paint0_linear_632_3973'
        gradientUnits='userSpaceOnUse'
      >
        <stop stop-color='#E9D7FE' />
        <stop offset='1' stop-color='#F1E6FF' />
      </linearGradient>
      <linearGradient
        x1='76.0004'
        y1='37.1539'
        x2='123.359'
        y2='46.5826'
        id='paint1_linear_632_3973'
        gradientUnits='userSpaceOnUse'
      >
        <stop stop-color='#E9D7FE' />
        <stop offset='1' stop-color='#F4EBFF' />
      </linearGradient>
    </defs>
  </svg>
);
