import { FC, useRef, PropsWithChildren } from 'react';

import { observer } from 'mobx-react-lite';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button.tsx';
import { Editor } from '@ui/form/Editor/Editor.tsx';
import { ParticipantsSelectGroup } from '@organization/components/Timeline/PastZone/events/email/compose-email/ParticipantsSelectGroup';
import { ModeChangeButtons } from '@organization/components/Timeline/PastZone/events/email/compose-email/EmailResponseModeChangeButtons';

export interface ComposeEmailProps extends PropsWithChildren {
  modal: boolean;
  replyToId?: string;
  onDiscard: () => void;
  emailUseCase: TimelineEmailUsecase;
  onModeChange?: (status: 'reply' | 'reply-all' | 'forward') => void;
}

export const ComposeEmail: FC<ComposeEmailProps> = observer(
  ({ onModeChange, onDiscard, modal, emailUseCase, replyToId }) => {
    const myRef = useRef<HTMLDivElement>(null);
    const height =
      modal && (myRef?.current?.getBoundingClientRect()?.height || 0) + 96;

    return (
      <form
        onSubmit={(e) => {
          e.preventDefault();
        }}
        className={cn(
          modal
            ? 'border-dashed border-t-[1px] border-gray-200 bg-grayBlue-50 rounded-none max-h-[50vh]'
            : 'bg-white rounded-lg max-h-[100%]',
          'rounded-b-2xl py-4 px-6 overflow-visible pt-1',
        )}
      >
        {!!onModeChange && (
          <div style={{ position: 'relative' }}>
            <ModeChangeButtons handleModeChange={onModeChange} />
          </div>
        )}
        <div ref={myRef}>
          <ParticipantsSelectGroup modal={modal} emailUseCase={emailUseCase} />
        </div>
        <div
          className='w-full mt-2'
          style={{
            maxHeight: modal ? `calc(50vh - ${height}px) !important` : 'auto',
          }}
        >
          <Editor
            showToolbarBottom
            dataTest='timeline-email-editor'
            placeholderClassName={'text-sm'}
            namespace='timeline-email-editor'
            placeholder={'Write something here...'}
            defaultHtmlValue={emailUseCase.emailContent}
            className='text-sm cursor-text email-editor h-full'
            onChange={(e) => {
              emailUseCase.updateEmailContent(e);
            }}
          >
            <div className='flex gap-2'>
              {onDiscard && (
                <Button size='xs' variant={'ghost'} onClick={onDiscard}>
                  Discard
                </Button>
              )}

              <Button
                size='xs'
                loadingText='Sending...'
                isLoading={emailUseCase.isSending}
                onClick={() => {
                  emailUseCase.createEmail(replyToId);
                }}
              >
                Send
              </Button>
            </div>
          </Editor>
        </div>
      </form>
    );
  },
);
