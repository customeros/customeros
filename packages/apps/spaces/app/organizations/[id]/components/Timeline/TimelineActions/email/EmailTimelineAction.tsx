import React, { useEffect, useRef } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { ComposeEmail } from '@organization/components/Timeline/events/email/compose-email/ComposeEmail';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useTimelineActionEmailContext } from '../TimelineActionsContext/TimelineActionEmailContext';

interface EmailTimelineActionProps {
  onScrollBottom: () => void;
}

export const EmailTimelineAction: React.FC<EmailTimelineActionProps> = ({
  onScrollBottom,
}) => {
  const {
    isEmailEditorOpen,
    remirrorProps,
    isSending,
    onCreateEmail,
    formId,
    state,
    handleExitEditorAndCleanData,
    showConfirmationDialog,
    closeConfirmationDialog,
  } = useTimelineActionEmailContext();
  const virtuoso = useRef(null);

  useEffect(() => {
    if (isEmailEditorOpen) {
      onScrollBottom();
    }
  }, [isEmailEditorOpen]);

  return (
    <>
      <Box
        bg={'#F9F9FB'}
        borderTop='1px dashed'
        borderTopColor='gray.200'
        pt={isEmailEditorOpen ? 6 : 0}
        pb={isEmailEditorOpen ? 2 : 8}
        mt={-4}
      >
        {isEmailEditorOpen && (
          <SlideFade in={true}>
            <Box
              ref={virtuoso}
              borderRadius={'md'}
              boxShadow={'lg'}
              m={6}
              mt={2}
              bg={'white'}
              border='1px solid var(--gray-100, #F2F4F7)'
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
      </Box>
      <ConfirmDeleteDialog
        label='Discard this email?'
        description='Saving draft emails is not possible at the moment. Would you like to continue to discard this email?'
        confirmButtonLabel='Discard email'
        isOpen={showConfirmationDialog}
        onClose={closeConfirmationDialog}
        onConfirm={() => handleExitEditorAndCleanData()}
        isLoading={false}
      />
    </>
  );
};
