import React, { useState } from 'react';
import { CardBody } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { EmailMetaDataEntry } from './EmailMetaDataEntry';
import { useTimelineEventPreviewContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { getEmailParticipantsByType } from '@organization/src/components/Timeline/events/email/utils';
import { ComposeEmail } from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail';
import { getEmailParticipantsNameAndEmail } from '@spaces/utils/getParticipantsName';
import Image from 'next/image';
import { useForm } from 'react-inverted-form';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail.dto';
import { handleSendEmail } from '@organization/src/components/Timeline/events/email/compose-email/utils';
import { useSearchParams } from 'next/navigation';
import { useSession } from 'next-auth/react';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useDisclosure } from '@ui/utils';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { htmlToProsemirrorNode } from 'remirror';
import { InteractionEvent } from '@graphql/types';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';
import { TimelinePreviewBackdrop } from '@organization/src/components/Timeline/preview/TimelinePreviewBackdrop';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';

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
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const event = modalContent as InteractionEvent;
  const subject = event?.interactionSession?.name || '';
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
    stringHandler: htmlToProsemirrorNode,
    content: '',
  });
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
    const params = new URLSearchParams(searchParams?.toString() ?? '');

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

  return (
    <TimelinePreviewBackdrop onCloseModal={handleClosePreview}>
      <Flex flexDir='column' maxH='calc(100vh - 5rem)' fontSize='sm'>
        <TimelineEventPreviewHeader
          //@ts-expect-error alias
          date={event.date}
          name={event.interactionSession?.name ?? ''}
          onClose={handleClosePreview}
          copyLabel='Copy link to this email'
        />

        <CardBody mt={0} p='6' pt='4' overflow='auto'>
          <Flex direction='row' justify='space-between' mb={3}>
            <Flex
              direction='column'
              align='flex-start'
              maxWidth='calc(100% - 70px)'
              overflow='hidden'
              textOverflow='ellipsis'
              fontSize='sm'
            >
              <EmailMetaDataEntry entryType='From' content={event?.sentBy} />
              <EmailMetaDataEntry entryType='To' content={to} />
              {!!cc.length && (
                <EmailMetaDataEntry entryType='CC' content={cc} />
              )}
              {!!bcc.length && (
                <EmailMetaDataEntry entryType='BCC' content={bcc} />
              )}
              <EmailMetaDataEntry entryType='Subject' content={subject} />
            </Flex>
            <div>
              <Image
                src={'/backgrounds/organization/post-stamp.webp'}
                alt='Email'
                width={48}
                height={70}
              />
            </div>
          </Flex>

          {event?.content && (
            <HtmlContentRenderer htmlContent={event.content} />
          )}
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
      </Flex>
    </TimelinePreviewBackdrop>
  );
};
