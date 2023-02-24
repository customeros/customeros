import { RefObject, useEffect, useRef } from 'react';

export const useDetectClickOutside = <T extends HTMLElement>(
  ref: RefObject<T>,
  callback: () => void,
) => {
  const callbackRef = useRef(callback);

  useEffect(() => {
    callbackRef.current = callback;
  }, [callback]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        callbackRef.current();
      }
    };

    // Timeout added to prevent automatically calling callback in case when eg. entering new mode
    setTimeout(
      () => document.addEventListener('click', handleClickOutside),
      100,
    );

    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  }, [ref]);
};
