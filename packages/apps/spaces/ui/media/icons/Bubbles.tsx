import { Icon, IconProps } from '@ui/media/Icon';

export const Bubbles = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='bubbles'>
      <circle
        id='Ellipse 23'
        cx='6.5'
        cy='7.5'
        r='2.5'
        stroke='currentColor'
        strokeWidth='2'
      />
      <circle
        id='Ellipse 25'
        cx='16.5'
        cy='5.5'
        r='1.5'
        fill='currentColor'
        stroke='currentColor'
        strokeWidth='2'
      />
      <circle
        id='Ellipse 24'
        cx='14.5'
        cy='15.5'
        r='4.5'
        stroke='currentColor'
        strokeWidth='2'
      />
    </g>
  </Icon>
);
