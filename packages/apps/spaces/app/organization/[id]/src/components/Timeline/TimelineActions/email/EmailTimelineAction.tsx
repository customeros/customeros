import React, { useEffect, useRef } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { ComposeEmail } from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail';
import { useTimelineActionEmailContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionEmailContext';
import { useTimelineActionContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionContext';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';

interface EmailTimelineActionProps {
  onScrollBottom: () => void;
}

export const EmailTimelineAction: React.FC<EmailTimelineActionProps> = ({
  onScrollBottom,
}) => {
  const {
    remirrorProps,
    isSending,
    onCreateEmail,
    formId,
    state,
    checkCanExitSafely,
  } = useTimelineActionEmailContext();
  const { openedEditor, showEditor } = useTimelineActionContext();
  const isEmail = openedEditor === 'email';
  const virtuoso = useRef(null);

  useEffect(() => {
    if (isEmail) {
      onScrollBottom();
    }
  }, [isEmail]);

  const handleClose = () => {
    const canClose = checkCanExitSafely();

    if (canClose) {
      showEditor(null);
    }
  };

  return (
    <>
      {isEmail && (
        <SlideFade in={true}>
          <Box
            ref={virtuoso}
            borderRadius={'md'}
            boxShadow={'lg'}
            m={6}
            mt={2}
            bg={'white'}
            border='1px solid'
            borderColor='gray.100'
          >
            <ComposeEmail
              formId={formId}
              modal={false}
              to={state.values.to}
              cc={state.values.cc}
              bcc={state.values.bcc}
              onSubmit={onCreateEmail}
              isSending={isSending}
              remirrorProps={remirrorProps}
            >
              <KeymapperClose onClose={handleClose} />
            </ComposeEmail>
          </Box>
        </SlideFade>
      )}
    </>
  );
};
