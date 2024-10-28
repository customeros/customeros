import { useParams } from 'react-router-dom';
import { useRef, useMemo, useState, useEffect } from 'react';

import { LexicalEditor } from 'lexical';
import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types';

import { cn } from '@ui/utils/cn';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { Modal, ModalPortal, ModalContent } from '@ui/overlay/Modal';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml';

import { useUndoRedo } from '../hooks';

interface EmailEditorModalProps {
  isEditorOpen: boolean;
  handleCancel: () => void;
  data: { subject: string; bodyTemplate: string; action: FlowActionType };
  handleEmailDataChange: (args: {
    subject: string;
    bodyTemplate: string;
  }) => void;
}

export const EmailEditorModal = observer(
  ({
    isEditorOpen,
    handleEmailDataChange,
    data,
    handleCancel,
  }: EmailEditorModalProps) => {
    const id = useParams().id as string;
    const inputRef = useRef<LexicalEditor>(null);

    const [subject, setSubject] = useState(data?.subject ?? '');
    const [bodyTemplate, setBodyTemplate] = useState(data?.bodyTemplate ?? '');
    const { takeSnapshot } = useUndoRedo();

    useEffect(() => {
      if (isEditorOpen) {
        setSubject(data?.subject ?? '');
        setBodyTemplate(data?.bodyTemplate ?? '');

        if (
          data.action !== FlowActionType.EMAIL_REPLY &&
          data?.subject?.trim()?.length === 0
        ) {
          setTimeout(() => {
            inputRef.current?.focus();
          }, 0);
        }
      }
    }, [isEditorOpen]);

    const store = useStore();

    const flow = store.flows.value.get(id)?.value?.name;
    const placeholder = useMemo(() => getRandomEmailPrompt(), [isEditorOpen]);

    const handleSave = () => {
      handleEmailDataChange({ subject, bodyTemplate });

      setTimeout(() => {
        takeSnapshot();
      }, 0);
    };

    const variables = store.flowEmailVariables?.value.get('CONTACT')?.variables;

    return (
      <Modal modal={false} open={isEditorOpen}>
        <ModalPortal>
          <ModalContent
            onKeyDown={(e) => e.stopPropagation()}
            className='w-full h-full flex flex-col align-middle items-center max-w-full top-0 cursor-default overflow-y-auto '
          >
            <div className='flex justify-between bg-white pt-4 pb-2 mb-[60px] w-[570px] sticky top-0 z-[50]'>
              <div className='flex items-center text-sm'>
                <span>{flow}</span>
                <ChevronRight className='size-3 mx-1 text-gray-400' />
                <span className='mr-2 cursor-default'>
                  {data.action === FlowActionType.EMAIL_NEW
                    ? 'Send Email'
                    : 'Reply to Email'}
                </span>
              </div>
              <div className='flex items-center gap-2'>
                <Button
                  size='xs'
                  variant='ghost'
                  onClick={() => {
                    setSubject(data.subject);
                    setBodyTemplate(data.bodyTemplate);
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
              <div>
                <Editor
                  size='md'
                  usePlainText
                  ref={inputRef}
                  placeholder='Subject'
                  namespace='flow-email-editor-subject'
                  onChange={(html) => setSubject(extractPlainText(html))}
                  defaultHtmlValue={convertPlainTextToHtml(subject ?? '')}
                  placeholderClassName='text-lg font-medium h-auto cursor-text'
                  className={cn(
                    `text-lg font-medium h-auto cursor-text email-editor-subject`,
                    {
                      'pointer-events-none text-gray-400':
                        data.action === FlowActionType.EMAIL_REPLY,
                    },
                  )}
                />
              </div>
              <div className='h-[60vh] mb-10'>
                <Editor
                  placeholder={placeholder}
                  variableOptions={variables}
                  dataTest='flow-email-editor'
                  namespace='flow-email-editor'
                  defaultHtmlValue={bodyTemplate}
                  onChange={(e) => setBodyTemplate(e)}
                  className='text-base cursor-text email-editor h-full'
                ></Editor>
              </div>
            </div>
          </ModalContent>
        </ModalPortal>
      </Modal>
    );
  },
);

const emailPrompts = [
  'Write something {{contact_first_name}} will want to share with their boss',
  'Craft an email that makes {{contact_first_name}} say "Wow!"',
  'Compose an email {{contact_first_name}} will quote in their presentation',
  "Make {{contact_first_name}} feel like they've discovered a hidden treasure",
  'Write an email that makes {{contact_first_name}} rethink their strategy',
  "Write something {{contact_first_name}} can't get from a Google search",
  'Compose the email that ends {{contact_first_name}}’s decision paralysis',
  "Write an email {{contact_first_name}} can't ignore",
  'Write something that makes {{contact_first_name}} feel stupid for not replying',
  'Write something that makes {{contact_first_name}} say, “Yes, this is what we need!”',
  'Show {{contact_first_name}} what they’re missing—start typing...',
  'Type an email that helps {{contact_first_name}} win',
  'Write something {{contact_first_name}} remember',
];

function getRandomEmailPrompt(): string {
  const randomIndex = Math.floor(Math.random() * emailPrompts.length);

  return emailPrompts[randomIndex];
}
