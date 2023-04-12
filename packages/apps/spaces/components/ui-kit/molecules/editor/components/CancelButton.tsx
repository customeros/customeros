import React from 'react';
import { Refresh } from '../../../atoms';
import { CommandButton, useRemirrorContext } from '@remirror/react';

export const CancelButton = () => {
  const { commands } = useRemirrorContext({
    autoUpdate: true,
  });
  const handleResetEditor = () => {
    commands.resetContent();
  };

  return (
    <CommandButton
      commandName='Reset editor'
      label='Cancel'
      onSelect={handleResetEditor}
      icon={<Refresh style={{ transform: 'scale(0.75)' }} />}
      enabled
      style={{
        maxHeight: '32px',
      }}
    />
  );
};
