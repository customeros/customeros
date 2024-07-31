import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Uz = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#1eb53a' d='M0 320h640v160H0z' />
    <path fill='#0099b5' d='M0 0h640v160H0z' />
    <path fill='#ce1126' d='M0 153.6h640v172.8H0z' />
    <path fill='#fff' d='M0 163.2h640v153.6H0z' />
    <circle r='57.6' cy='76.8' cx='134.4' fill='#fff' />
    <circle r='57.6' cy='76.8' cx='153.6' fill='#0099b5' />
    <g fill='#fff' transform='translate(261.1 122.9)scale(1.92)'>
      <g id='uz-e'>
        <g id='uz-d'>
          <g id='uz-c'>
            <g id='uz-b'>
              <path id='uz-a' d='M0-6-1.9-.3 1 .7' />
              <use
                width='100%'
                height='100%'
                xlinkHref='#uz-a'
                transform='scale(-1 1)'
              />
            </g>
            <use
              width='100%'
              height='100%'
              xlinkHref='#uz-b'
              transform='rotate(72)'
            />
          </g>
          <use
            width='100%'
            height='100%'
            xlinkHref='#uz-b'
            transform='rotate(-72)'
          />
          <use
            width='100%'
            height='100%'
            xlinkHref='#uz-c'
            transform='rotate(144)'
          />
        </g>
        <use y='-24' width='100%' height='100%' xlinkHref='#uz-d' />
        <use y='-48' width='100%' height='100%' xlinkHref='#uz-d' />
      </g>
      <use x='24' width='100%' height='100%' xlinkHref='#uz-e' />
      <use x='48' width='100%' height='100%' xlinkHref='#uz-e' />
      <use x='-48' width='100%' height='100%' xlinkHref='#uz-d' />
      <use x='-24' width='100%' height='100%' xlinkHref='#uz-d' />
      <use x='-24' y='-24' width='100%' height='100%' xlinkHref='#uz-d' />
    </g>
  </svg>
);
