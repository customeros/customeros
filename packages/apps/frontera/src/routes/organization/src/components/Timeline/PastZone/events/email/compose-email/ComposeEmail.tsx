import { useField } from 'react-inverted-form';
import { FC, useRef, PropsWithChildren } from 'react';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button.tsx';
import { Editor } from '@ui/form/Editor/Editor.tsx';
import { ParticipantsSelectGroup } from '@organization/components/Timeline/PastZone/events/email/compose-email/ParticipantsSelectGroup';
import { ModeChangeButtons } from '@organization/components/Timeline/PastZone/events/email/compose-email/EmailResponseModeChangeButtons';

export interface ComposeEmailProps extends PropsWithChildren {
  formId: string;
  modal: boolean;
  isSending: boolean;
  onSubmit: () => void;
  attendees: Array<string>;
  to: Array<{ label: string; value: string }>;
  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;
  onModeChange?: (status: 'reply' | 'reply-all' | 'forward') => void;
}

export const ComposeEmail: FC<ComposeEmailProps> = ({
  onModeChange,
  formId,
  modal,
  isSending,
  onSubmit,
  attendees,
  to,
  cc,
  bcc,
}) => {
  const myRef = useRef<HTMLDivElement>(null);
  const height =
    modal && (myRef?.current?.getBoundingClientRect()?.height || 0) + 96;
  const {
    getInputProps,
    state: { value },
  } = useField('content', formId);

  const { onChange } = getInputProps();

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
        <ParticipantsSelectGroup
          to={to}
          cc={cc}
          bcc={bcc}
          modal={modal}
          formId={formId}
          attendees={attendees}
        />
      </div>
      <div
        className='w-full mt-2'
        style={{
          maxHeight: modal ? `calc(50vh - ${height}px) !important` : 'auto',
        }}
      >
        <Editor
          placeholder={''}
          showToolbarBottom
          defaultHtmlValue={value}
          onChange={(e) => onChange(e)}
          dataTest='timeline-email-editor'
          namespace='timeline-email-editor'
          className='text-base cursor-text email-editor h-full'
        >
          <Button
            size='xs'
            onClick={onSubmit}
            isLoading={isSending}
            loadingText='Sending...'
          >
            Send
          </Button>
        </Editor>
      </div>
    </form>
  );
};
