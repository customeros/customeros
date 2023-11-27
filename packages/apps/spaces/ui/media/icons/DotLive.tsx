import { Icon, IconProps } from '@ui/media/Icon';

export const DotLive = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <circle cx='12' cy='12' r='8.25' stroke='currentColor' strokeWidth='0.5' />
    <circle cx='12' cy='12' r='4' fill='currentColor' />
  </Icon>
);
