import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Outlook = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 32 28'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g id='outlook'>
      <g id='rectangle'>
        <rect x='10' width='20' height='28' rx='2' fill='#1066B5' />
        <rect
          x='10'
          width='20'
          height='28'
          rx='2'
          fill='url(#paint0_linear_1334_1221)'
        />
      </g>
      <rect
        id='rectangle_2'
        x='10'
        y='3'
        width='10'
        height='10'
        fill='#32A9E7'
      />
      <rect
        id='rectangle_3'
        x='10'
        y='13'
        width='10'
        height='10'
        fill='#167EB4'
      />
      <rect
        id='rectangle_4'
        x='20'
        y='13'
        width='10'
        height='10'
        fill='#32A9E7'
      />
      <rect
        id='rectangle_5'
        x='20'
        y='3'
        width='10'
        height='10'
        fill='#58D9FD'
      />
      <g id='Mask Group'>
        <mask
          id='mask0_1334_1221'
          maskUnits='userSpaceOnUse'
          x='8'
          y='12'
          width='24'
          height='16'
        >
          <path
            id='Rectangle 3'
            d='M8 12H30C31.1046 12 32 12.8954 32 14V26C32 27.1046 31.1046 28 30 28H10C8.89543 28 8 27.1046 8 26V12Z'
            fill='url(#paint1_linear_1334_1221)'
          />
        </mask>
        <g mask='url(#mask0_1334_1221)'>
          <path id='vector' d='M32 12V16H30V12H32Z' fill='#135298' />
          <path
            id='vector_2'
            d='M32 28V14L7 28H32Z'
            fill='url(#paint2_linear_1334_1221)'
          />
          <path
            id='vector_3'
            d='M8 28V14L33 28H8Z'
            fill='url(#paint3_linear_1334_1221)'
          />
        </g>
      </g>
      <path
        id='rectangle_6'
        d='M8 10C8 8.34315 9.34315 7 11 7H17C18.6569 7 20 8.34315 20 10V22C20 23.6569 18.6569 25 17 25H8V10Z'
        fill='black'
        fill-opacity='0.3'
      />
      <rect
        id='rectangle_7'
        y='5'
        width='18'
        height='18'
        rx='2'
        fill='url(#paint4_linear_1334_1221)'
      />
      <path
        id='O'
        d='M14 14.0693V13.903C14 11.0222 11.9272 9 9.01582 9C6.08861 9 4 11.036 4 13.9307V14.097C4 16.9778 6.07278 19 9 19C11.9114 19 14 16.964 14 14.0693ZM11.6424 14.097C11.6424 16.0083 10.5665 17.1579 9.01582 17.1579C7.46519 17.1579 6.37342 15.9806 6.37342 14.0693V13.903C6.37342 11.9917 7.44937 10.8421 9 10.8421C10.5348 10.8421 11.6424 12.0194 11.6424 13.9307V14.097Z'
        fill='white'
      />
    </g>
    <defs>
      <linearGradient
        id='paint0_linear_1334_1221'
        x1='10'
        y1='14'
        x2='30'
        y2='14'
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#064484' />
        <stop offset='1' stopColor='#0F65B5' />
      </linearGradient>
      <linearGradient
        id='paint1_linear_1334_1221'
        x1='8'
        y1='24.7692'
        x2='32'
        y2='24.7692'
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#1B366F' />
        <stop offset='1' stopColor='#2657B0' />
      </linearGradient>
      <linearGradient
        id='paint2_linear_1334_1221'
        x1='32'
        y1='21'
        x2='8'
        y2='21'
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#44DCFD' />
        <stop offset='0.453125' stopColor='#259ED0' />
      </linearGradient>
      <linearGradient
        id='paint3_linear_1334_1221'
        x1='8'
        y1='21'
        x2='32'
        y2='21'
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#259ED0' />
        <stop offset='1' stopColor='#44DCFD' />
      </linearGradient>
      <linearGradient
        id='paint4_linear_1334_1221'
        x1='0'
        y1='14'
        x2='18'
        y2='14'
        gradientUnits='userSpaceOnUse'
      >
        <stop stopColor='#064484' />
        <stop offset='1' stopColor='#0F65B5' />
      </linearGradient>
    </defs>
  </svg>
);
