import Lottie from 'react-lottie';
import React, { useRef, useState } from 'react';

import { Button } from '@ui/form/Button/Button';
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
      size='md'
      variant='ghost'
      onClick={() => triggerConfetti()}
      className='w-full px-3 text-gray-500 font-normal text-sm'
      leftIcon={<ArrowsRight className='size-5 text-gray-500' />}
    >
      <div className='flex justify-between flex-1 items-center'>
        <span>Up next</span>
        {!!unseenCount &&
          (unseenCount >= 99 ? (
            <div
              onClick={triggerConfetti}
              className='relative overflow-visible w-[35px] h-[33px]'
            >
              {isBursted ? (
                <>
                  <Lottie
                    width={100}
                    height={100}
                    options={defaultOptions}
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
                  <p className='text-white absolute text-xs z-[1] left-[18px] top-[7px]'>
                    99+
                  </p>
                </>
              )}
            </div>
          ) : (
            <div className='size-5 flex items-center justify-center rounded-xl border border-gray-300 font-normal'>
              {unseenCount}
            </div>
          ))}
      </div>
      <audio ref={audioRef} src='/soundEffects/99_audio.mp4' />
    </Button>
  );
};
