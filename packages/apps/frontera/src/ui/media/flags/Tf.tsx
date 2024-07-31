import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Tf = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <path id='tf-a' fill='#fff' d='m0-21 12.3 38L-20-6.5h40L-12.3 17z' />
    </defs>
    <path fill='#002395' d='M0 0h640v480H0z' />
    <path fill='#fff' d='M0 0h292.8v196.8H0z' />
    <path fill='#002395' d='M0 0h96v192H0z' />
    <path fill='#ed2939' d='M192 0h96v192h-96z' />
    <path
      fill='#fff'
      d='m426 219.6 15.4 24.6h44V330l-33-51.6-44.4 70.8h21.6l22.8-40.8 46.8 84 46.8-84 22.8 40.8h21.6L546 278.4 513 330v-47.4h19.8l14.7-23.4H513v-15h44l15.4-24.6zm51.6 105h-48v16.8h48zm91.2 0h-48v16.8h48z'
    />
    <use
      x='416'
      y='362'
      width='100%'
      height='100%'
      xlinkHref='#tf-a'
      transform='scale(1.2)'
    />
    <use
      x='371'
      y='328'
      width='100%'
      height='100%'
      xlinkHref='#tf-a'
      transform='scale(1.2)'
    />
    <use
      x='461'
      y='328'
      width='100%'
      height='100%'
      xlinkHref='#tf-a'
      transform='scale(1.2)'
    />
    <use
      x='333'
      y='227'
      width='100%'
      height='100%'
      xlinkHref='#tf-a'
      transform='scale(1.2)'
    />
    <use
      x='499'
      y='227'
      width='100%'
      height='100%'
      xlinkHref='#tf-a'
      transform='scale(1.2)'
    />
  </svg>
);
