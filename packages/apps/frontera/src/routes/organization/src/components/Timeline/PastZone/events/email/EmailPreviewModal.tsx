import { useMemo, useState } from 'react';
import { useForm } from 'react-inverted-form';
import { VirtuosoHandle } from 'react-virtuoso';
import { useSearchParams } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { EmailParticipant, InteractionEvent } from '@graphql/types';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { getEmailParticipantsNameAndEmail } from '@utils/getParticipantsName';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';
import { getEmailParticipantsByType } from '@organization/components/Timeline/PastZone/events/email/utils';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/PastZone/hooks/updateCacheWithNewEvent';
import { TimelinePreviewBackdrop } from '@organization/components/Timeline/shared/TimelineEventPreview/TimelinePreviewBackdrop';
import { ComposeEmailContainer } from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmailContainer';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmail.dto';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { EmailMetaDataEntry } from './EmailMetaDataEntry';

import postStamp from '/backgrounds/organization/post-stamp.webp';

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';
declare type FieldProps = {
  error?: string;
  meta: {
    pristine: boolean;
    hasError: boolean;
    isTouched: boolean;
  };
};

declare type Fields<T> = Record<keyof T, FieldProps>;

const checkPristine = (
  fieldsData: Partial<Fields<ComposeEmailDtoI>>,
): boolean => {
  return Object.values(fieldsData).every((e) => e.meta.pristine);
};

const checkEmpty = (values: Partial<ComposeEmailDtoI>): boolean => {
  return (
    !values.from || !values.fromProvider || !values.to || values.to.length === 0
  );
};
const formId = 'compose-email-preview-modal';
interface EmailPreviewModalProps {
  invalidateQuery: () => void;
  virtuosoRef?: React.RefObject<VirtuosoHandle>;
}

export const EmailPreviewModal = ({
  invalidateQuery,
  virtuosoRef,
}: EmailPreviewModalProps) => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const { open: isOpen, onOpen, onClose } = useDisclosure();

  const event = modalContent as InteractionEvent;
  const subject = event?.interactionSession?.name || '';

  const updateTimelineCache = useUpdateCacheWithNewEvent(virtuosoRef);
  const [searchParams] = useSearchParams();
  const store = useStore();
  const [mode, setMode] = useState(REPLY_MODE);
  const [isSending, setIsSending] = useState(false);
  const { to, cc, bcc } = getEmailParticipantsByType(event?.sentTo || []);

  const fromParticipants = event?.sentBy || [];
  const from = getEmailParticipantsNameAndEmail(fromParticipants, 'value');
  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    from: from,
    to: getEmailParticipantsNameAndEmail(to, 'value'),
    cc: getEmailParticipantsNameAndEmail(cc, 'value'),
    bcc: getEmailParticipantsNameAndEmail(bcc, 'value'),
    subject: `Re: ${subject}`,
    content: '',
  });
  const [timelineMeta] = useTimelineMeta();
  const queryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );

  const { state, setDefaultValues } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,
  });

  //will remain here until we figure out if we need it or not
  // const handleResetEditor = () => {
  //   const context = remirrorProps.getContext();
  //   if (context) {
  //     context.commands.resetContent();
  //   }
  // };

  const handleEmailSendSuccess = async (response: unknown) => {
    await updateTimelineCache(response, queryKey);

    setDefaultValues(defaultValues);
    // no timeout needed is this case as the event id is created when this is called
    invalidateQuery();
    setIsSending(false);
    closeModal();
  };

  const handleEmailSendError = () => {
    setIsSending(false);
  };

  const handleModeChange = (newMode: string) => {
    let newDefaultValues = defaultValues;

    function removeDuplicates(
      emailTO: Array<{ label: string; [x: string]: string }>,
      emailCC: Array<{ label: string; [x: string]: string }>,
    ): Array<{ label: string; [x: string]: string }> {
      const uniqueValuesSet = new Set(emailTO.map((email) => email.value));

      const filteredCC = emailCC.filter(
        (email) => !uniqueValuesSet.has(email.value),
      );

      return filteredCC;
    }

    const newTo = from[0].value.includes(store.session.value.profile.email)
      ? [
          ...getEmailParticipantsNameAndEmail(
            [
              ...to.filter(
                (e) =>
                  e.emailParticipant.email !==
                  store.session.value.profile.email,
              ),
            ],
            'value',
          ),
        ]
      : from;

    if (newMode === REPLY_MODE) {
      newDefaultValues = new ComposeEmailDto({
        from: '',
        fromProvider: '',
        to: newTo,
        cc: [],
        bcc: [],
        subject: `Re: ${subject}`,
        content: mode === FORWARD_MODE ? '' : state.values.content,
      });
    }

    if (newMode === REPLY_ALL_MODE) {
      const newCC = [
        ...getEmailParticipantsNameAndEmail(
          [
            ...cc,
            ...to.filter(
              (e) =>
                e.emailParticipant.email !== store.session.value.profile.email,
            ),
          ],
          'value',
        ),
      ];
      const newBCC = getEmailParticipantsNameAndEmail(bcc, 'value');

      newDefaultValues = new ComposeEmailDto({
        from: '',
        fromProvider: '',
        to: [...newTo],
        cc: removeDuplicates(newTo, newCC),
        bcc: newBCC,
        subject: `Re: ${subject}`,
        content: mode === FORWARD_MODE ? '' : state.values.content,
      });
    }

    if (newMode === FORWARD_MODE) {
      newDefaultValues = new ComposeEmailDto({
        from: '',
        fromProvider: '',
        to: [],
        cc: [],
        bcc: [],
        subject: `Re: ${subject}`,
        content: `${state.values.content}${event.content}`,
      });
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
      ((!areParticipantFieldsPristine && !areFieldsEmpty) ||
        (!subjectField.meta.pristine && !subject.length)) &&
      !isFormEmpty;

    if (isFormPristine || !showConfirmationDialog) {
      handleExitEditorAndCleanData();
    } else {
      onOpen();
    }
  };

  const handleSubmit = () => {
    const from = state.values.from?.value ?? '';
    const fromProvider = state.values.from?.provider ?? '';
    const to = [...state.values.to].map(({ value }) => value);
    const cc = [...state.values.cc].map(({ value }) => value);
    const bcc = [...state.values.bcc].map(({ value }) => value);
    const params = new URLSearchParams(searchParams?.toString() ?? '');

    setIsSending(true);

    const id = params.get('events') ?? undefined;

    store.mail.send(
      {
        from,
        fromProvider,
        to,
        cc,
        bcc,
        replyTo: id,
        content: state.values.content,
        subject: state.values.subject,
      },
      {
        onSuccess: handleEmailSendSuccess,
        onError: handleEmailSendError,
      },
    );
  };

  const filteredParticipants = useMemo(
    () => ({
      attendees: (event?.interactionSession?.attendedBy
        ?.map((i) => i as EmailParticipant)
        .map((i) => i.emailParticipant.email) ?? []) as string[],
      to: state.values.to?.filter((e) => !!e.value || !!e?.label),
      cc: state.values.cc?.filter((e) => !!e.value || !!e?.label),
      bcc: state.values.bcc?.filter((e) => !!e.value || !!e?.label),
    }),
    [event, state.values.to, state.values.cc, state.values.bcc],
  );

  return (
    <TimelinePreviewBackdrop onCloseModal={handleClosePreview}>
      <div className='flex flex-col max-h-[calc(100vh-5rem)] text-sm max-w-[700px]'>
        <TimelineEventPreviewHeader
          //@ts-expect-error alias
          date={event.date}
          onClose={handleClosePreview}
          copyLabel='Copy link to this email'
          name={event.interactionSession?.name ?? ''}
        />

        <div className='mt-0 p-6 pt-4 overflow-auto'>
          <div className='flex flex-row justify-between mb-3'>
            <div className='flex flex-col items-start max-w-[calc(100%-70px)] overflow-hidden text-sm line-clamp-1'>
              <EmailMetaDataEntry entryType='From' content={event?.sentBy} />
              <EmailMetaDataEntry content={to} entryType='To' />
              {!!cc.length && (
                <EmailMetaDataEntry content={cc} entryType='CC' />
              )}
              {!!bcc.length && (
                <EmailMetaDataEntry content={bcc} entryType='BCC' />
              )}
              <EmailMetaDataEntry content={subject} entryType='Subject' />
            </div>
            <div>
              <img alt='Email' src={postStamp} className='w-[48px] h-[70px]' />
            </div>
          </div>

          {event?.content && (
            <HtmlContentRenderer htmlContent={event.content} />
          )}
        </div>

        <ComposeEmailContainer
          {...filteredParticipants}
          modal
          formId={formId}
          isSending={isSending}
          onSubmit={handleSubmit}
          onClose={handleClosePreview}
          onModeChange={handleModeChange}
        />

        <ConfirmDeleteDialog
          isOpen={isOpen}
          isLoading={false}
          colorScheme='primary'
          onConfirm={handleSubmit}
          confirmButtonLabel='Send'
          label={`Send this email?`}
          cancelButtonLabel='Discard'
          onClose={handleExitEditorAndCleanData}
          description={`You have typed an unsent email. Do you want to send it, or discard it?`}
        />
      </div>
    </TimelinePreviewBackdrop>
  );
};
