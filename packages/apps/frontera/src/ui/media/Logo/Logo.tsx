import { SVGAttributes } from 'react';

import { twMerge } from 'tailwind-merge';

export type LogoName =
  | 'google'
  | 'facebook'
  | 'linkedin'
  | 'github'
  | 'twitter'
  | 'youtube'
  | 'instagram'
  | 'default'
  | 'grain'
  | 'fathom';
interface LogoProps extends SVGAttributes<SVGElement> {
  name: LogoName;
  className?: string;
}

export const Logo = ({
  name,
  width,
  height,
  stroke,
  viewBox,
  className,
  strokeWidth,
  ...props
}: LogoProps) => (
  <svg
    fillRule='evenodd'
    clipRule='evenodd'
    width={width ?? 24}
    height={height ?? 24}
    viewBox={viewBox ?? '0 0 24 24'}
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <use xlinkHref={`/logos.svg#${name}`} />
  </svg>
);
