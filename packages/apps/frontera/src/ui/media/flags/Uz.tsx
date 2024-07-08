import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Uz = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#1eb53a' d='M0 320h640v160H0z' />
    <path fill='#0099b5' d='M0 0h640v160H0z' />
    <path fill='#ce1126' d='M0 153.6h640v172.8H0z' />
    <path fill='#fff' d='M0 163.2h640v153.6H0z' />
    <circle cx='134.4' cy='76.8' r='57.6' fill='#fff' />
    <circle cx='153.6' cy='76.8' r='57.6' fill='#0099b5' />
    <g fill='#fff' transform='translate(261.1 122.9)scale(1.92)'>
      <g id='uz-e'>
        <g id='uz-d'>
          <g id='uz-c'>
            <g id='uz-b'>
              <path id='uz-a' d='M0-6-1.9-.3 1 .7' />
              <use
                xlinkHref='#uz-a'
                width='100%'
                height='100%'
                transform='scale(-1 1)'
              />
            </g>
            <use
              xlinkHref='#uz-b'
              width='100%'
              height='100%'
              transform='rotate(72)'
            />
          </g>
          <use
            xlinkHref='#uz-b'
            width='100%'
            height='100%'
            transform='rotate(-72)'
          />
          <use
            xlinkHref='#uz-c'
            width='100%'
            height='100%'
            transform='rotate(144)'
          />
        </g>
        <use xlinkHref='#uz-d' width='100%' height='100%' y='-24' />
        <use xlinkHref='#uz-d' width='100%' height='100%' y='-48' />
      </g>
      <use xlinkHref='#uz-e' width='100%' height='100%' x='24' />
      <use xlinkHref='#uz-e' width='100%' height='100%' x='48' />
      <use xlinkHref='#uz-d' width='100%' height='100%' x='-48' />
      <use xlinkHref='#uz-d' width='100%' height='100%' x='-24' />
      <use xlinkHref='#uz-d' width='100%' height='100%' x='-24' y='-24' />
    </g>
  </svg>
);
