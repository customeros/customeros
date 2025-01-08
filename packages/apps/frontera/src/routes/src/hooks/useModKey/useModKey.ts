import { useEffect } from 'react';

export const useModKey = (
  key: string,
  callback: (e: KeyboardEvent) => void,
  options: { when?: boolean; targetRef?: React.RefObject<HTMLElement> } = {
    when: true,
  },
) => {
  useEffect(() => {
    if (!options.when) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === key && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        callback(e);
      }
    };

    const target = options.targetRef?.current || document;

    target.addEventListener('keydown', handleKeyDown as EventListener);

    return () => {
      target.removeEventListener('keydown', handleKeyDown as EventListener);
    };
  }, [options.when, callback, options.targetRef]);
};
