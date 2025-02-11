import { useRef } from 'react';

import { LexicalEditor } from 'lexical';
import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types';

import { cn } from '@ui/utils/cn';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { EmailVariableName } from '@graphql/types';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { Modal, ModalPortal, ModalContent } from '@ui/overlay/Modal';
import { extractPlainText } from '@ui/form/Editor/utils/extractPlainText';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml';

interface EmailEditorModalProps {
  subject: string;
  flowName?: string;
  placeholder: string;
  bodyTemplate: string;
  isReadOnly?: boolean;
  isEditorOpen: boolean;
  handleSave: () => void;
  action?: FlowActionType;
  handleCancel: () => void;
  variables?: Array<EmailVariableName>;
  setSubject: (subject: string) => void;
  setBodyTemplate: (bodyTemplate: string) => void;
}

export const EmailEditorModal = observer(
  ({
    isEditorOpen,
    subject,
    bodyTemplate,
    handleCancel,
    setSubject,
    setBodyTemplate,
    flowName,
    placeholder,
    variables,
    action,
    handleSave,
    isReadOnly,
  }: EmailEditorModalProps) => {
    const inputRef = useRef<LexicalEditor>(null);
    const editorRef = useRef<LexicalEditor>(null);

    return (
      <Modal modal={false} open={isEditorOpen}>
        <ModalPortal>
          <ModalContent
            onKeyDown={(e) => e.stopPropagation()}
            className='w-full h-full flex flex-col align-middle items-center max-w-full top-0 cursor-default overflow-y-auto'
          >
            <div className='flex justify-between bg-white pt-4 pb-2 mb-[60px] w-[570px] sticky top-0 z-[50]'>
              <div className='flex items-center text-sm'>
                <span>{flowName}</span>
                <ChevronRight className='size-3 mx-1 text-gray-400' />
                <span className='mr-2 cursor-default'>
                  {action === FlowActionType.EMAIL_NEW
                    ? 'Send Email'
                    : 'Reply to Email'}
                </span>
              </div>
              <div className='flex items-center gap-2'>
                <Button
                  size='xs'
                  variant='ghost'
                  onClick={() => {
                    setSubject(subject);
                    setBodyTemplate(bodyTemplate);
                    handleCancel();
                  }}
                >
                  Discard
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
                  isReadOnly={isReadOnly}
                  variableOptions={variables}
                  namespace='flow-email-editor-subject'
                  onChange={(html) => setSubject(extractPlainText(html))}
                  defaultHtmlValue={convertPlainTextToHtml(subject ?? '')}
                  placeholderClassName='text-lg font-medium h-auto cursor-text'
                  onKeyDown={(e) => {
                    if (e.key === 'Tab') {
                      e.preventDefault();
                      editorRef.current?.focus();
                    }
                  }}
                  className={cn(
                    `text-lg font-medium h-auto cursor-text email-editor-subject`,
                    {
                      'pointer-events-none text-gray-400':
                        action === FlowActionType.EMAIL_REPLY,
                    },
                  )}
                />
              </div>
              <div className='h-[60vh] mb-10'>
                <Editor
                  ref={editorRef}
                  isReadOnly={isReadOnly}
                  placeholder={placeholder}
                  variableOptions={variables}
                  dataTest='flow-email-editor'
                  namespace='flow-email-editor'
                  defaultHtmlValue={bodyTemplate}
                  onChange={(e) => setBodyTemplate(e)}
                  className='text-base cursor-text email-editor h-full'
                />
              </div>
            </div>
          </ModalContent>
        </ModalPortal>
      </Modal>
    );
  },
);
