import { Icon, IconProps } from '@ui/media/Icon';

export const ArrowsRight = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <path
      d='M4 7H15M15 7L11 11M15 7L11 3M4 17H20M20 17L16 21M20 17L16 13'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </Icon>
);
