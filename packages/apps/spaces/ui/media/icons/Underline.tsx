import { Icon, IconProps } from '@ui/media/Icon';

export const Underline = (props: IconProps) => (
  <Icon viewBox='0 0 16 16' fill='none' boxSize='4' {...props}>
    <path
      d='M12 2.66667V7.33334C12 9.54248 10.2091 11.3333 8 11.3333C5.79086 11.3333 4 9.54248 4 7.33334V2.66667M2.66667 14H13.3333'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </Icon>
);
