import { Icon, IconProps } from '@ui/media/Icon';

export const ChevronDown = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <path
      d='M6 9L12 15L18 9'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </Icon>
);
