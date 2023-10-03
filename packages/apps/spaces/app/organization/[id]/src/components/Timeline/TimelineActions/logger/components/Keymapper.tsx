import { useEffect } from 'react';
import { useKeymap, useCommands } from '@remirror/react';

import { useTimelineActionLogEntryContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionLogEntryContext';

export const Keymapper = () => {
  const { focus } = useCommands();
  const { onCreateLogEntry } = useTimelineActionLogEntryContext();

  useKeymap('Mod-Enter', ({ next }) => {
    onCreateLogEntry();
    return next();
  });

  useEffect(() => {
    focus('start');
  }, []);

  return null;
};
