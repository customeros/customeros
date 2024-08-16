import { useEffect } from 'react';

export const useModKey = (
  key: string,
  callback: (e: KeyboardEvent) => void,
  options: { when?: boolean } = { when: true },
) => {
  useEffect(() => {
    if (!options.when) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === key && e.metaKey) {
        e.preventDefault();
        callback(e);
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [options.when, callback]);
};
