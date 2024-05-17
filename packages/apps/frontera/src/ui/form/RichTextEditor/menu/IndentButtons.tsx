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
        onClick={() => {
          commands.decreaseIndent();
        }}
        icon={<RightIndent className='text-inherit' />}
      />
      <ToolbarButton
        label='Outdent'
        onClick={() => {
          commands.increaseIndent();
        }}
        icon={<LeftIndent className='text-inherit' />}
      />
    </div>
  );
};
