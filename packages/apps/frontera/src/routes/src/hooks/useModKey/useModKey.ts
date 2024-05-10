import { useEffect } from 'react';

export const useModKey = (
  key: string,
  callback: (e: KeyboardEvent) => void,
) => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === key && e.metaKey) {
        callback(e);
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, []);
};
