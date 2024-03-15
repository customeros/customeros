import { useRef } from 'react';

import { useOutsideClick } from '@ui/utils';
import { Flex, FlexProps } from '@ui/layout/Flex';
import { pulseOpacity } from '@ui/utils/keyframes';

interface ReminderPostitProps extends FlexProps {
  isMutating?: boolean;
  onClickOutside?: (e: Event) => void;
}

export const ReminderPostit = ({
  children,
  isMutating,
  onClickOutside = () => undefined,
  ...rest
}: ReminderPostitProps) => {
  const ref = useRef(null);

  useOutsideClick({ ref, handler: onClickOutside });

  return (
    <Flex
      ref={ref}
      position='relative'
      w='321px'
      m='6'
      mt='2'
      pointerEvents={isMutating ? 'none' : 'auto'}
      animation={
        isMutating ? `${pulseOpacity} 0.7s alternate ease-in-out` : undefined
      }
      {...rest}
    >
      <Flex
        h='10px'
        w='calc(100% - 5px)'
        bottom='-5px'
        position='absolute'
        filter={'blur(5px)'}
        bg='rgba(0, 0, 0, 0.15)'
        transform={'rotate(2deg)'}
      />
      <Flex w='full' zIndex={1} boxShadow='md' bg='yellow.100' flexDir='column'>
        <Flex
          h='28px'
          w='full'
          bgGradient='linear(to-tl, rgba(196, 196, 196, 0.00) 66.96%, rgba(0, 0, 0, 0.06) 99.59%)'
        />
        {children}
      </Flex>
    </Flex>
  );
};
