import { useEffect } from 'react';
import { useKeymap, useCommands } from '@remirror/react';

export const KeymapperCreate = ({ onCreate }: { onCreate: () => void }) => {
  const { focus } = useCommands();

  useKeymap('Mod-Enter', ({ next }) => {
    onCreate();
    return next();
  });

  useEffect(() => {
    focus('start');
  }, []);

  return null;
};
