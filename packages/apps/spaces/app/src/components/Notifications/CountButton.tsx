import Lottie from 'react-lottie';
import React, { useRef, useState } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Badge } from '@ui/presentation/Badge';
import { RedBalloon } from '@ui/media/icons/RedBalloon';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';

import animationData from './confetti.json';

interface CountBalloonProps {
  unseenCount?: number;
}

export const CountButton: React.FC<CountBalloonProps> = ({ unseenCount }) => {
  const [isBursted, setIsBursted] = useState(false);
  const audioRef = useRef<HTMLAudioElement>(null);
  const defaultOptions = {
    loop: false,
    autoplay: true,
    animationData: animationData,
    rendererSettings: {
      preserveAspectRatio: 'xMidYMid slice',
    },
  };
  const triggerConfetti = () => {
    if (unseenCount && unseenCount >= 99) {
      setIsBursted(true);
      audioRef?.current?.play();
    }
  };

  return (
    <Button
      px='3'
      w='full'
      size='md'
      variant='ghost'
      fontSize='sm'
      textDecoration='none'
      fontWeight='regular'
      justifyContent='flex-start'
      borderRadius='md'
      color={'gray.500'}
      onClick={() => triggerConfetti()}
      leftIcon={<ArrowsRight className='size-5' />}
      _focus={{
        boxShadow: 'sidenavItemFocus',
      }}
    >
      <Flex justifyContent='space-between' flex={1} alignItems='center'>
        <span>Up next</span>
        {!!unseenCount &&
          (unseenCount >= 99 ? (
            <Box
              position='relative'
              overflow='visible'
              w={35}
              h={33}
              onClick={triggerConfetti}
            >
              {isBursted ? (
                <>
                  <Lottie
                    options={defaultOptions}
                    height={100}
                    width={100}
                    style={{
                      position: 'absolute',
                      zIndex: '10',
                      top: '-32px',
                      left: '-26px',
                    }}
                  />
                </>
              ) : (
                <>
                  <RedBalloon className='size-[53px] absolute left-0 z-10' />
                  <Text
                    color='white'
                    position='absolute'
                    zIndex={1}
                    fontSize='xs'
                    left='18px'
                    top='7px'
                  >
                    99+
                  </Text>
                </>
              )}
            </Box>
          ) : (
            <Badge
              w={5}
              h={5}
              display='flex'
              alignItems='center'
              justifyContent='center'
              variant='outline'
              borderRadius='xl'
              boxShadow='none'
              border='1px solid'
              borderColor='gray.300'
              fontWeight='regular'
            >
              {unseenCount}
            </Badge>
          ))}
      </Flex>
      <audio ref={audioRef} src='/soundEffects/99_audio.mp4' />
    </Button>
  );
};
