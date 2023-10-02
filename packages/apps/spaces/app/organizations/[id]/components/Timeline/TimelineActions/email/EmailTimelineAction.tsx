import React, { useEffect, useRef } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { ComposeEmail } from '@organization/components/Timeline/events/email/compose-email/ComposeEmail';
import { useTimelineActionEmailContext } from '@organization/components/Timeline/TimelineActions/context/TimelineActionEmailContext';
import { useTimelineActionContext } from '@organization/components/Timeline/TimelineActions/context/TimelineActionContext';

interface EmailTimelineActionProps {
  onScrollBottom: () => void;
}

export const EmailTimelineAction: React.FC<EmailTimelineActionProps> = ({
  onScrollBottom,
}) => {
  const { remirrorProps, isSending, onCreateEmail, formId, state } =
    useTimelineActionEmailContext();
  const { openedEditor } = useTimelineActionContext();
  const isEmail = openedEditor === 'email';
  const virtuoso = useRef(null);

  useEffect(() => {
    if (isEmail) {
      onScrollBottom();
    }
  }, [isEmail]);

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
            />
          </Box>
        </SlideFade>
      )}
    </>
  );
};
