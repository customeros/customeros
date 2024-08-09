import { useRef, useEffect, useCallback } from 'react';

type KeyboardEventKey = KeyboardEvent['key'];

interface SequentialShortcutOptions {
  when?: boolean;
}

export function useSequentialShortcut(
  firstKey: KeyboardEventKey,
  secondKey: KeyboardEventKey,
  callback: () => void,
  options: SequentialShortcutOptions = { when: true },
): void {
  const isFirstKeyPressedRef = useRef<boolean>(false);
  const timerRef = useRef<number | null>(null);

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (event.key.toLowerCase() === firstKey.toLowerCase()) {
        isFirstKeyPressedRef.current = true;

        if (timerRef.current !== null) {
          clearTimeout(timerRef.current);
        }

        timerRef.current = window.setTimeout(() => {
          isFirstKeyPressedRef.current = false;
        }, 1000);
      } else if (
        isFirstKeyPressedRef.current &&
        event.key.toLowerCase() === secondKey.toLowerCase()
      ) {
        event.preventDefault();
        callback();
        isFirstKeyPressedRef.current = false;

        if (timerRef.current !== null) {
          clearTimeout(timerRef.current);
        }
      } else {
        isFirstKeyPressedRef.current = false;

        if (timerRef.current !== null) {
          clearTimeout(timerRef.current);
        }
      }
    },
    [firstKey, secondKey, callback],
  );

  useEffect(() => {
    if (!options.when) return;

    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);

      if (timerRef.current !== null) {
        clearTimeout(timerRef.current);
      }
    };
  }, [handleKeyDown, options.when]);
}
