import { FC, useRef, PropsWithChildren } from 'react';

import { observer } from 'mobx-react-lite';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { Button } from '@ui/form/Button/Button.tsx';
import { Editor } from '@ui/form/Editor/Editor.tsx';
import { ParticipantsSelectGroup } from '@organization/components/Timeline/PastZone/events/email/compose-email/ParticipantsSelectGroup';

export interface ComposeEmailProps extends PropsWithChildren {
  modal: boolean;
  replyToId?: string;
  onDiscard: () => void;
  emailUseCase: TimelineEmailUsecase;
}

export const ComposeEmail: FC<ComposeEmailProps> = observer(
  ({ onDiscard, modal, emailUseCase, replyToId }) => {
    const myRef = useRef<HTMLDivElement>(null);
    const height =
      modal && (myRef?.current?.getBoundingClientRect()?.height || 0) + 96;

    return (
      <form
        onSubmit={(e) => {
          e.preventDefault();
        }}
      >
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
                <Button
                  size='xs'
                  variant={'ghost'}
                  onClick={onDiscard}
                  dataTest='timeline-email-discard'
                >
                  Discard
                </Button>
              )}

              <Button
                size='xs'
                loadingText='Sending...'
                dataTest='timeline-email-send'
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
