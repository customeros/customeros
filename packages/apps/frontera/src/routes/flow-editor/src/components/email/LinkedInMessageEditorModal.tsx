import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types';

import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { Modal, ModalPortal, ModalContent } from '@ui/overlay/Modal';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml';

import { useUndoRedo } from '../../hooks';

interface EmailEditorModalProps {
  isEditorOpen: boolean;
  handleCancel: () => void;
  data: { action: FlowActionType; messageTemplate: string };
  handleEmailDataChange: (args: { messageTemplate: string }) => void;
}

export const LinkedInMessageEditorModal = observer(
  ({
    isEditorOpen,
    handleEmailDataChange,
    data,
    handleCancel,
  }: EmailEditorModalProps) => {
    const id = useParams().id as string;
    const store = useStore();

    const [messageTemplate, setMessageTemplate] = useState(
      data?.messageTemplate ?? '',
    );
    const { takeSnapshot } = useUndoRedo();

    useEffect(() => {
      if (isEditorOpen) {
        setMessageTemplate(data?.messageTemplate ?? '');
      }
    }, [isEditorOpen, data]);

    const flow = store.flows.value.get(id)?.value?.name;
    const variables = store.flowEmailVariables?.value.get('CONTACT')?.variables;

    const handleSave = async () => {
      try {
        handleEmailDataChange({ messageTemplate });
        setTimeout(() => {
          takeSnapshot();
        }, 0);
      } catch (error) {
        console.error('Error saving email:', error);
        store.ui.toastError('Error saving email', 'email-save-error');
      }
    };

    return (
      <Modal modal={false} open={isEditorOpen}>
        <ModalPortal>
          <ModalContent
            onKeyDown={(e) => e.stopPropagation()}
            className='w-full h-full flex flex-col align-middle items-center max-w-full top-0 cursor-default overflow-y-auto'
          >
            <div className='flex justify-between bg-white pt-4 pb-2 mb-[60px] w-[570px] sticky top-0 z-[50]'>
              <div className='flex items-center text-sm'>
                <span>{flow}</span>
                <ChevronRight className='size-3 mx-1 text-gray-400' />
                <span className='mr-2 cursor-default'>
                  {data?.action === FlowActionType.LINKEDIN_MESSAGE
                    ? 'Send LinkedIn Message'
                    : 'Send Connection Request'}
                </span>
              </div>
              <div className='flex items-center gap-2'>
                <Button
                  size='xs'
                  variant='ghost'
                  onClick={() => {
                    setMessageTemplate(data.messageTemplate);
                    handleCancel();
                  }}
                >
                  Cancel
                </Button>
                <Button
                  size='xs'
                  variant='outline'
                  leftIcon={<Check />}
                  onClick={handleSave}
                >
                  Done
                </Button>
              </div>
            </div>
            <div className='w-[570px] relative'>
              <div className='h-[60vh] mb-10'>
                <Editor
                  size='md'
                  usePlainText
                  className='cursor-text'
                  variableOptions={variables}
                  namespace='opportunity-next-step'
                  placeholderClassName='cursor-text'
                  onChange={(html) =>
                    setMessageTemplate(extractPlainText(html))
                  }
                  defaultHtmlValue={convertPlainTextToHtml(
                    messageTemplate ?? '',
                  )}
                />
              </div>
            </div>
          </ModalContent>
        </ModalPortal>
      </Modal>
    );
  },
);
