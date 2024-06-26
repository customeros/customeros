import { FC, useRef, PropsWithChildren } from 'react';

import { cn } from '@ui/utils/cn';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { BasicEditorToolbar } from '@ui/form/RichTextEditor/menu/BasicEditorToolbar';
import {
  RemirrorProps,
  BasicEditorExtentions,
} from '@ui/form/RichTextEditor/types';
import { KeymapperCreate } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperCreate';
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
  remirrorProps: RemirrorProps<BasicEditorExtentions>;
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
  remirrorProps,
  children,
}) => {
  const myRef = useRef<HTMLDivElement>(null);
  const height =
    modal && (myRef?.current?.getBoundingClientRect()?.height || 0) + 96;

  return (
    <form
      className={cn(
        modal
          ? 'border-dashed border-t-[1px] border-gray-200 bg-grayBlue-50 rounded-none max-h-[50vh]'
          : 'bg-white rounded-lg max-h-[100%]',
        'rounded-b-2xl py-4 px-6 overflow-visible pt-1',
      )}
      onSubmit={(e) => {
        e.preventDefault();
      }}
    >
      {!!onModeChange && (
        <div style={{ position: 'relative' }}>
          <ModeChangeButtons handleModeChange={onModeChange} />
        </div>
      )}
      <div ref={myRef}>
        <ParticipantsSelectGroup
          attendees={attendees}
          to={to}
          cc={cc}
          bcc={bcc}
          modal={modal}
          formId={formId}
        />
      </div>
      <div
        className='w-full'
        style={{
          maxHeight: modal ? `calc(50vh - ${height}px) !important` : 'auto',
        }}
      >
        <RichTextEditor
          {...remirrorProps}
          formId={formId}
          name='content'
          showToolbar
        >
          {children}
          <KeymapperCreate onCreate={onSubmit} />
          <BasicEditorToolbar isSending={isSending} onSubmit={onSubmit} />
        </RichTextEditor>
      </div>
    </form>
  );
};
