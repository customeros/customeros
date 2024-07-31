import { useCommands } from '@remirror/react';

import { LeftIndent } from '@ui/media/icons/LeftIndent';
import { RightIndent } from '@ui/media/icons/RightIndent';

import { ToolbarButton } from './ToolbarButton';

export const IndentButtons = () => {
  const commands = useCommands();

  return (
    <div className='flex gap-2'>
      <ToolbarButton
        label='Indent'
        icon={<RightIndent className='text-inherit' />}
        onClick={() => {
          commands.decreaseIndent();
        }}
      />
      <ToolbarButton
        label='Outdent'
        icon={<LeftIndent className='text-inherit' />}
        onClick={() => {
          commands.increaseIndent();
        }}
      />
    </div>
  );
};
