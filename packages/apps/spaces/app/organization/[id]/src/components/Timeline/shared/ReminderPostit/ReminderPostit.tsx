import { useRef, useMemo, useState } from 'react';

import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
import { Flex, FlexProps } from '@ui/layout/Flex';
import { pulseOpacity } from '@ui/utils/keyframes';

interface ReminderPostitProps extends FlexProps {
  owner?: string;
  isFocused?: boolean;
  isMutating?: boolean;
  onClickOutside?: (e: Event) => void;
}

const rotations = ['rotate(2deg)', 'rotate(-2deg)', 'rotate(0deg)'];
const rgadients = [
  'linear(to-tl, rgba(196, 196, 196, 0.00) 20%, rgba(0, 0, 0, 0.06) 100%)',
  'linear(to-tr, rgba(196, 196, 196, 0.00) 20%, rgba(0, 0, 0, 0.06) 100%)',
  'linear(to-bl, rgba(196, 196, 196, 0.00) 20%, rgba(0, 0, 0, 0.03) 100%)',
];

const getRandomStyles = () => {
  const index = Math.floor(Math.random() * rotations.length);

  return [rotations[index], rgadients[index]];
};

export const ReminderPostit = ({
  owner,
  children,
  isFocused,
  isMutating,
  onClickOutside = () => undefined,
  ...rest
}: ReminderPostitProps) => {
  const ref = useRef(null);
  const [isHovered, setIsHovered] = useState(false);
  const [rotation, gradient] = useMemo(() => getRandomStyles(), []);

  useOutsideClick({ ref, handler: onClickOutside });

  return (
    <Flex
      ref={ref}
      position='relative'
      w='321px'
      m='6'
      mt='2'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      pointerEvents={isMutating ? 'none' : 'auto'}
      animation={
        isMutating ? `${pulseOpacity} 0.7s alternate ease-in-out` : undefined
      }
      {...rest}
    >
      <Flex
        h='calc(100% - 28px)'
        w='calc(100% - 10px)'
        bottom='-4px'
        left='5px'
        position='absolute'
        filter={isFocused || isHovered ? 'blur(7px)' : 'blur(3px)'}
        bg={
          isFocused || isHovered ? 'rgba(0, 0, 0, 0.2)' : 'rgba(0, 0, 0, 0.07)'
        }
        transition='all 0.1s ease-in-out'
        transform={isFocused || isHovered ? 'unset' : rotation}
      />
      <Flex w='full' zIndex={1} bg='yellow.100' flexDir='column'>
        <Flex h='24px' w='full' align='center' bgGradient={gradient}>
          {owner && (
            <Text pt='3' fontSize='xs' color='gray.500' pl='4' fontWeight='400'>
              {owner} added
            </Text>
          )}
        </Flex>
        {children}
      </Flex>
    </Flex>
  );
};
