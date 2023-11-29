import { Icon, IconProps } from '@ui/media/Icon';

export const DotLive = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <circle
      cx='12'
      cy='12'
      r='8.25'
      fill='#DCFAE6'
      stroke='#ABEFC6'
      strokeWidth='0.5'
    />
    <circle cx='12' cy='12' r='4' fill='#17B26A' />
  </Icon>
);
