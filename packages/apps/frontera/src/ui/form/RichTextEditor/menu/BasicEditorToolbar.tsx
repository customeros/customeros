import { FC } from 'react';

import { useActive, useCommands } from '@remirror/react';

import { Quote } from '@ui/media/icons/Quote';
import { Button } from '@ui/form/Button/Button';

import { ListButtons } from './ListButtons';
import { ToolbarButton } from './ToolbarButton';
import { TextFormatButtons } from './TextFormatButtons';

export const BasicEditorToolbar: FC<{
  isSending: boolean;
  onSubmit: () => void;
}> = ({ isSending, onSubmit }) => {
  const { toggleBlockquote, focus } = useCommands();
  const active = useActive();

  return (
    <div className='flex justify-between items-center flex-1 min-h-8'>
      <div className='flex bg-transparent w-full'>
        <TextFormatButtons />
        <div className='h-8 bg-gray-200 w-[1px] mr-[2px]' />
        <ListButtons />
        <div className='h-8 bg-gray-200 w-[1px] mr-[2px]' />
        <ToolbarButton
          label='Quote'
          onClick={() => {
            toggleBlockquote();
            focus();
          }}
          isActive={active.blockquote()}
          icon={<Quote className='text-inherit' />}
        />
      </div>
      <Button
        className='customeros-remirror-submit-button font-semibold rounded-lg px-3 py-1 text-sm'
        variant='outline'
        colorScheme='gray'
        size='sm'
        isDisabled={isSending}
        isLoading={isSending}
        loadingText='Sending'
        onClick={onSubmit}
      >
        Send
      </Button>
    </div>
  );
};
