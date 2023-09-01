import React, { useState, useRef } from 'react';
import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { EmailMetaDataEntry } from './EmailMetaDataEntry';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { useCopyToClipboard } from '@spaces/hooks/useCopyToClipboard';
import sanitizeHtml from 'sanitize-html';
import { DateTimeUtils } from '@spaces/utils/date';
import { getEmailParticipantsByType } from '@organization/components/Timeline/events/email/utils';
import CopyLink from '@spaces/atoms/icons/CopyLink';
import Times from '@spaces/atoms/icons/Times';
import { ComposeEmail } from '@organization/components/Timeline/events/email/compose-email/ComposeEmail';
import { getEmailParticipantsNameAndEmail } from '@spaces/utils/getParticipantsName';
import Image from 'next/image';
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
import { htmlToProsemirrorNode } from 'remirror';
import { InteractionEvent } from '@graphql/types';
import { RichTextPreview } from '@ui/form/RichTextEditor/RichTextPreview';
import { useOutsideClick } from '@ui/utils';

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';
declare type FieldProps = {
  meta: {
    pristine: boolean;
    hasError: boolean;
    isTouched: boolean;
  };
  error?: string;
};

declare type Fields<T> = Record<keyof T, FieldProps>;
const checkPristine = (
  fieldsData: Partial<Fields<ComposeEmailDtoI>>,
): boolean => {
  return Object.values(fieldsData).every((e) => e.meta.pristine);
};

const checkEmpty = (values: Partial<ComposeEmailDtoI>): boolean => {
  return Object.values(values).every((e) => !e.length);
};

interface EmailPreviewModalProps {
  invalidateQuery: () => void;
}

export const EmailPreviewModal: React.FC<EmailPreviewModalProps> = ({
  invalidateQuery,
}) => {
  const { closeModal, isModalOpen, modalContent } =
    useTimelineEventPreviewContext();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cardRef = useRef<HTMLDivElement>(null);

  const event = modalContent as InteractionEvent;

  const subject = event?.interactionSession?.name || '';
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
    stringHandler: htmlToProsemirrorNode,
    content: '',
  });
  const [_, copy] = useCopyToClipboard();
  const searchParams = useSearchParams();
  const { data: session } = useSession();
  const [mode, setMode] = useState(REPLY_MODE);
  const [isSending, setIsSending] = useState(false);
  const { to, cc, bcc } = getEmailParticipantsByType(event?.sentTo || []);
  const from = getEmailParticipantsNameAndEmail(event?.sentBy || [], 'value');
  const formId = 'compose-email-preview-modal';

  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: getEmailParticipantsNameAndEmail(
      [...(event?.sentBy ?? []), ...(to ?? [])],
      'value',
    ),
    cc: getEmailParticipantsNameAndEmail(cc, 'value'),
    bcc: getEmailParticipantsNameAndEmail(bcc, 'value'),
    subject: `Re: ${subject}`,
    content: '',
  });

  const { state, setDefaultValues } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,

    stateReducer: (state, action, next) => {
      return next;
    },
  });

  const handleResetEditor = () => {
    const context = remirrorProps.getContext();
    if (context) {
      context.commands.resetContent();
    }
  };

  const handleEmailSendSuccess = () => {
    invalidateQuery();
    setIsSending(false);
    setDefaultValues(defaultValues);
    closeModal();
  };
  const handleEmailSendError = () => {
    setIsSending(false);
  };

  const handleModeChange = (newMode: string) => {
    let newDefaultValues = defaultValues;
    if (newMode === REPLY_MODE) {
      newDefaultValues = new ComposeEmailDto({
        to: from,
        cc: [],
        bcc: [],
        subject: `Re: ${subject}`,
        content: mode === FORWARD_MODE ? '' : state.values.content,
      });
    }
    if (newMode === REPLY_ALL_MODE) {
      const newTo = getEmailParticipantsNameAndEmail(to, 'value');
      const newCC = getEmailParticipantsNameAndEmail(cc, 'value');
      const newBCC = getEmailParticipantsNameAndEmail(bcc, 'value');
      newDefaultValues = new ComposeEmailDto({
        to: [...from, ...newTo],
        cc: newCC,
        bcc: newBCC,
        subject: `Re: ${subject}`,
        content: mode === FORWARD_MODE ? '' : state.values.content,
      });
      handleResetEditor();
    }
    if (newMode === FORWARD_MODE) {
      newDefaultValues = new ComposeEmailDto({
        to: [],
        cc: [],
        bcc: [],
        subject: `Re: ${subject}`,
        content: `${state.values.content}${event.content}`,
      });
      const prosemirrorNodeValue = htmlToProsemirrorNode({
        schema: remirrorProps.state.schema,
        content: `<p>${state.values.content} ${event.content}</p>`,
      });
      remirrorProps.getContext()?.setContent(prosemirrorNodeValue);
    }
    setMode(newMode);
    setDefaultValues(newDefaultValues);
  };

  const handleExitEditorAndCleanData = () => {
    setDefaultValues(defaultValues);
    onClose();
    closeModal();
  };

  const handleClosePreview = (): void => {
    const { content, subject, ...values } = state.values;
    const {
      content: contentField,
      subject: subjectField,
      ...fields
    } = state.fields;

    const isFormPristine = checkPristine(state.fields);
    const areParticipantFieldsPristine = checkPristine(fields);

    const isFormEmpty = !content.length || content === `<p style=""></p>`;
    const areFieldsEmpty = checkEmpty(values);

    const showConfirmationDialog =
      (!areParticipantFieldsPristine && !areFieldsEmpty) ||
      (!subjectField.meta.pristine && !subject.length) ||
      !isFormEmpty;

    if (isFormPristine || !showConfirmationDialog) {
      handleExitEditorAndCleanData();
    } else {
      onOpen();
    }
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

  useOutsideClick({
    ref: cardRef,
    handler: handleClosePreview,
  });

  if (!isModalOpen || !modalContent) {
    return null;
  }

  return (
    <div ref={cardRef}>
      <CardHeader
        pb={1}
        position='sticky'
        background='white'
        top={0}
        borderRadius='xl'
        onClick={(e) => e.stopPropagation()}
      >
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <div>
            <Heading size='sm' mb={2}>
              {event.interactionSession?.name}
            </Heading>
            <Text size='2xs' color='gray.500' fontSize='12px'>
              {DateTimeUtils.format(
                // @ts-expect-error this is correct (alias)
                event.date,
                DateTimeUtils.dateWithHour,
              )}
            </Text>
          </div>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link to this email' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this email'
                color='gray.500'
                size='sm'
                mr={1}
                icon={<CopyLink color='gray.500' height='18px' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                size='sm'
                icon={<Times color='gray.500' height='24px' />}
                onClick={handleClosePreview}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>

      <CardBody mt={0} maxHeight='50%' overflow='auto' pb={6}>
        <Flex direction='row' justify='space-between' mb={3}>
          <Flex
            direction='column'
            align='flex-start'
            maxWidth='calc(100% - 70px)'
            overflow='hidden'
            textOverflow='ellipsis'
          >
            <EmailMetaDataEntry entryType='From' content={event?.sentBy} />
            <EmailMetaDataEntry entryType='To' content={to} />
            {!!cc.length && <EmailMetaDataEntry entryType='CC' content={cc} />}
            {!!bcc.length && (
              <EmailMetaDataEntry entryType='BCC' content={bcc} />
            )}
            <EmailMetaDataEntry entryType='Subject' content={subject} />
          </Flex>
          <div>
            <Image
              src={'/backgrounds/organization/post-stamp.webp'}
              alt='Email'
              width={54}
              height={70}
              style={{
                filter: 'drop-shadow(0px 0.5px 1px #D8D8D8)',
              }}
            />
          </div>
        </Flex>

        <Text color='gray.700' size='sm'>
          {event?.content && (
            <RichTextPreview
              htmlContent={sanitizeHtml(event.content)}
              extensions={basicEditorExtensions}
            />
          )}
        </Text>
      </CardBody>
      <ComposeEmail
        formId={formId}
        onModeChange={handleModeChange}
        modal
        to={state.values.to}
        cc={state.values.cc}
        bcc={state.values.bcc}
        onSubmit={handleSubmit}
        isSending={isSending}
        remirrorProps={remirrorProps}
      />
      <ConfirmDeleteDialog
        label='Discard this email?'
        description='Saving draft emails is not possible at the moment. Would you like to continue to discard this email?'
        confirmButtonLabel='Discard email'
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={handleExitEditorAndCleanData}
        isLoading={false}
      />
    </div>
  );
};
