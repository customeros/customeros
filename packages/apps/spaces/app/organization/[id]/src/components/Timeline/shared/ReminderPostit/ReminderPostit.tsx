import { Flex } from '@ui/layout/Flex';

interface ReminderPostitProps {
  children: React.ReactNode;
}

export const ReminderPostit = ({ children }: ReminderPostitProps) => {
  return (
    <Flex position='relative' w='321px' m='6' mt='2'>
      <Flex
        h='10px'
        w='calc(100% - 5px)'
        bottom='-5px'
        position='absolute'
        filter={'blur(5px)'}
        bg='rgba(0, 0, 0, 0.15)'
        transform={'rotate(2deg)'}
      />
      <Flex flexDir='column' bg='yellow.100' w='full' zIndex={1} boxShadow='md'>
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
