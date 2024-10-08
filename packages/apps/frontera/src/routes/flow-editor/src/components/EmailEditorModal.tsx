import { useParams } from 'react-router-dom';
import React, { useRef, useMemo, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types.ts';

import { Input } from '@ui/form/Input';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

import { useUndoRedo } from '../hooks';

interface EmailEditorModalProps {
  isEditorOpen: boolean;
  data: { subject: string; bodyTemplate: string; action: FlowActionType };
  handleEmailDataChange: (args: {
    subject: string;
    bodyTemplate: string;
  }) => void;
}

export const EmailEditorModal = observer(
  ({ isEditorOpen, handleEmailDataChange, data }: EmailEditorModalProps) => {
    const id = useParams().id as string;
    const inputRef = useRef<HTMLInputElement>(null);

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

    return (
      <Modal open={isEditorOpen}>
        <ModalPortal>
          <ModalOverlay className='z-50'>
            <ModalContent
              onKeyDown={(e) => e.stopPropagation()}
              className='w-full h-full flex justify-center max-w-full top-0 cursor-default overflow-y-auto'
            >
              <div className='w-[570px]'>
                <div className='flex justify-between mt-4 mb-[68px]'>
                  <div className='flex items-center text-sm'>
                    <span>{flow}</span>
                    <ChevronRight className='size-3 mx-1 text-gray-400' />
                    <span className='mr-2 cursor-default'>
                      {data.action === FlowActionType.EMAIL_NEW
                        ? 'Send Email'
                        : 'Reply to Email'}
                    </span>
                  </div>
                  <Button
                    size='xs'
                    variant='ghost'
                    leftIcon={<Check />}
                    onClick={handleSave}
                  >
                    Done
                  </Button>
                </div>

                <Input
                  ref={inputRef}
                  value={subject}
                  variant='unstyled'
                  placeholder='Subject'
                  className='font-medium text-lg min-h-[auto]'
                  onChange={(e) => setSubject(e.target.value)}
                  disabled={data.action === FlowActionType.EMAIL_REPLY}
                />

                <Editor
                  placeholder={placeholder}
                  dataTest='flow-email-editor'
                  namespace='flow-email-editor'
                  defaultHtmlValue={bodyTemplate}
                  onChange={(e) => setBodyTemplate(e)}
                  className='mb-10 text-base cursor-text email-editor'
                ></Editor>
              </div>
            </ModalContent>
          </ModalOverlay>
        </ModalPortal>
      </Modal>
    );
  },
);

const emailPrompts = [
  "Write something they'll want to share with their boss",
  'Craft an email that makes them say "Wow!"',
  "Compose an email they'll quote in their presentation",
  "Make them feel like they've discovered hidden treasure",
  'Write an email that makes them rethink their strategy',
  "Write something they can't get from a Google search",
  'Compose the email that ends their decision paralysis',
  "Write an email they can't ignore",
  'Turn this blank canvas into a sales masterpiece',
  'Write something that makes them feel stupid for not replying',
  'Write something that makes them say, "Yes, this is what we need!"',
  "Show them what they're missing—start typing...",
  'Type an email that helps them win',
  "Write something they'll remember",
  'Make your email impossible to ignore',
  'Start an email that stands out',
];

function getRandomEmailPrompt(): string {
  const randomIndex = Math.floor(Math.random() * emailPrompts.length);

  return emailPrompts[randomIndex];
}
