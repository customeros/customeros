import { useRef } from 'react';

import { useOutsideClick } from '@ui/utils';
import { Flex, FlexProps } from '@ui/layout/Flex';

export const ReminderPostit = ({
  children,
  onClickOutside = () => undefined,
  ...rest
}: FlexProps & { onClickOutside?: (e: Event) => void }) => {
  const ref = useRef(null);

  useOutsideClick({ ref, handler: onClickOutside });

  return (
    <Flex ref={ref} position='relative' w='321px' m='6' mt='2' {...rest}>
      <Flex
        h='10px'
        w='calc(100% - 5px)'
        bottom='-5px'
        position='absolute'
        filter={'blur(5px)'}
        bg='rgba(0, 0, 0, 0.15)'
        transform={'rotate(2deg)'}
      />
      <Flex
        w='full'
        zIndex={1}
        boxShadow='md'
        bg='yellow.100'
        id='sticky-body'
        flexDir='column'
      >
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
