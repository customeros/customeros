import { useEffect } from 'react';
import { useKeymap, useCommands } from '@remirror/react';

export const KeymapperClose = ({ onClose }: { onClose: () => void }) => {
  const { focus } = useCommands();

  useKeymap('Escape', ({ next }) => {
    onClose();

    return next();
  });

  useEffect(() => {
    focus('start');
  }, []);

  return null;
};
