import { useAnimate } from 'framer-motion';
import { useEffect } from 'react';

export function useTagButtonSlideAnimation(isOpen: boolean) {
  const [scope, animate] = useAnimate();

  useEffect(() => {
    animate([
      [
        'div > div > div > span',
        isOpen
          ? {
              opacity: 1,
              scale: 1,
              transform: 'translateX(0px)',
              filter: 'blur(0px)',
            }
          : {
              opacity: 0,
              scale: 0.3,
              filter: 'blur(2px)',
              transform: 'translateX(100px)',
            },
        {
          duration: 0.2,
          delay: 0,
        },
      ],
    ]);
  }, [isOpen]);

  return scope;
}
