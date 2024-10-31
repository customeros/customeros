import { useEffect } from 'react';

import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { ComposeEmailContainer } from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmailContainer';
import { useTimelineActionContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';
import { useTimelineActionEmailContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionEmailContext';

export const EmailTimelineAction = () => {
  const { isSending, onCreateEmail, formId, state, checkCanExitSafely } =
    useTimelineActionEmailContext();
  const { virtuosoRef } = useTimelineRefContext();
  const { closeEditor } = useTimelineActionContext();

  useEffect(() => {
    virtuosoRef?.current?.scrollBy({ top: 300 });
  }, [virtuosoRef]);

  const handleClose = () => {
    const canClose = checkCanExitSafely();

    if (canClose) {
      closeEditor();
    }
  };

  return (
    <div className='rounded-md shadow-lg m-6 mt-2 bg-white border border-gray-100 max-w-[800px]'>
      <ComposeEmailContainer
        modal={false}
        attendees={[]}
        formId={formId}
        to={state.values.to}
        cc={state.values.cc}
        onClose={handleClose}
        isSending={isSending}
        bcc={state.values.bcc}
        onSubmit={onCreateEmail}
      >
        <KeymapperClose onClose={handleClose} />
      </ComposeEmailContainer>
    </div>
  );
};
