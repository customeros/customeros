import React, { useEffect, useRef, useState } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { Button } from '@ui/form/Button';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { ComposeEmail } from '@organization/components/Timeline/events/email/compose-email/ComposeEmail';
import { Icons } from '@ui/media/Icon';
import { useForm } from 'react-inverted-form';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/events/email/compose-email/ComposeEmail.dto';
import { handleSendEmail } from '@organization/components/Timeline/events/email/compose-email/utils';
import { useSearchParams } from 'next/navigation';
import { useSession } from 'next-auth/react';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useDisclosure } from '@ui/utils';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { Logger } from '@organization/components/Timeline/TimelineActions/logger/Logger';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';

interface TimelineActionsProps {
  onScrollBottom: () => void;
  invalidateQuery: () => void;
}

export const TimelineActions: React.FC<TimelineActionsProps> = ({
  onScrollBottom,
  invalidateQuery,
}) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [showEmailEditor, setShowEmailEditor] = React.useState(false);
  const [showLogger, setShowLogger] = React.useState(false);
  const [isSending, setIsSending] = useState(false);
  const virtuoso = useRef(null);
  const searchParams = useSearchParams();
  const { data: session } = useSession();
  const formId = 'compose-email-timeline-footer';
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: [],
    cc: [],
    bcc: [],
    subject: '',
    content: '',
  });
  useEffect(() => {
    if (showEmailEditor || showLogger) {
      onScrollBottom();
    }
  }, [showEmailEditor, showLogger]);

  const { state, reset } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,

    stateReducer: (state, action, next) => {
      return next;
    },
  });

  const handleEmailSendSuccess = () => {
    invalidateQuery();
    setIsSending(false);
    reset();
    setShowEmailEditor(false);
  };
  const handleEmailSendError = () => {
    setIsSending(false);
  };

  const handleSubmit = () => {
    const to = [...state.values.to].map(({ value }) => value);
    const cc = [...state.values.cc].map(({ value }) => value);
    const bcc = [...state.values.bcc].map(({ value }) => value);
    const params = new URLSearchParams(searchParams ?? '');

    setIsSending(true);
    const id = params.get('events');
    return handleSendEmail(
      state.values.content,
      to,
      cc,
      bcc,
      id,
      state.values.subject,
      handleEmailSendSuccess,
      handleEmailSendError,
      session?.user?.email,
    );
  };

  const handleExitEditorAndCleanData = () => {
    reset();
    onClose();
    setShowEmailEditor(false);
  };
  const handleCloseEmailEditor = () => {
    const isFormPristine = Object.values(state.fields)?.every(
      (e) => e.meta.pristine,
    );
    const isFormEmpty = Object.values(state.values)?.every((e) => !e.length);

    const showEmailEditorConfirmationDialog = !isFormPristine && !isFormEmpty;
    if (showEmailEditorConfirmationDialog) {
      onOpen();
    } else {
      handleExitEditorAndCleanData();
    }
  };

  const handleToggleEmailEditor = () => {
    if (!showEmailEditor) {
      setShowLogger(false);
      setShowEmailEditor(true);
    } else {
      handleCloseEmailEditor();
    }
  };

  const handleToggleLogger = () => {
    setShowEmailEditor(false);
    setShowLogger((prevState) => !prevState);
  };

  return (
    <Box bg='gray.25'>
      <ButtonGroup
        position='sticky'
        border='1px solid'
        borderColor='gray.200'
        p='2'
        borderRadius='full'
        bg='white'
        top='0'
        left='6'
        zIndex='1'
        transform='translateY(5px)'
      >
        <Button
          variant='outline'
          onClick={() => handleToggleEmailEditor()}
          borderRadius='3xl'
          size='xs'
          leftIcon={<Icons.Mail1 boxSize='4' />}
        >
          Email
        </Button>
        <Button
          variant='outline'
          onClick={handleToggleLogger}
          borderRadius='3xl'
          size='xs'
          leftIcon={<MessageChatSquare />}
        >
          Log
        </Button>
      </ButtonGroup>
      <Box
        bg={'#F9F9FB'}
        borderTop='1px dashed'
        borderTopColor='gray.200'
        pt={showEmailEditor || showLogger ? 6 : 0}
        pb={showEmailEditor || showLogger ? 2 : 8}
        mt={-4}
      >
        {showEmailEditor && (
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
                onSubmit={handleSubmit}
                isSending={isSending}
                remirrorProps={remirrorProps}
              />
            </Box>
          </SlideFade>
        )}
        {showLogger && (
          <SlideFade in={true}>
            <Box
              ref={virtuoso}
              borderRadius={'md'}
              boxShadow={'lg'}
              m={6}
              mt={2}
              p={6}
              bg={'white'}
              border='1px solid var(--gray-100, #F2F4F7)'
            >
              <Logger />
            </Box>
          </SlideFade>
        )}
      </Box>
      <ConfirmDeleteDialog
        label='Discard this email?'
        description='Saving draft emails is not possible at the moment. Would you like to continue to discard this email?'
        confirmButtonLabel='Discard email'
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={() => handleExitEditorAndCleanData()}
        isLoading={false}
      />
    </Box>
  );
};
