import { useMemo, useState } from 'react';
import { useForm } from 'react-inverted-form';
import { VirtuosoHandle } from 'react-virtuoso';
import { useParams, useSearchParams } from 'react-router-dom';

import { useKey } from 'rooks';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { useStore } from '@shared/hooks/useStore';
import { InteractionEvent } from '@graphql/types';
import { useChannel } from '@shared/hooks/useChannel';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
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
  const orgId = useParams().id as string;

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

  const { currentUserId } = useChannel(`finder:${store.session.value.tenant}`);

  const emailUseCase = useMemo(
    () => new TimelineEmailUsecase(orgId, [], currentUserId),
    [orgId, event.id, isOpen],
  );

  const { state, setDefaultValues } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,
  });

  const handleModeChange = (newMode: string) => {
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

    const reSubject = subject.toLowerCase().includes('re:')
      ? subject
      : `Re: ${subject}`;

    if (newMode === REPLY_MODE) {
      emailUseCase.toSelector.select(newTo);
      emailUseCase.ccSelector.reset();
      emailUseCase.bccSelector.reset();
      emailUseCase.updateSubject(reSubject);
      emailUseCase.updateEmailContent(
        mode === FORWARD_MODE ? '' : emailUseCase.emailContent,
      );
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

      emailUseCase.toSelector.select(newTo);
      emailUseCase.ccSelector.select(removeDuplicates(newTo, newCC));
      emailUseCase.bccSelector.select(newBCC);
      emailUseCase.updateSubject(reSubject);
      emailUseCase.updateEmailContent(
        mode === FORWARD_MODE ? '' : emailUseCase.emailContent,
      );
    }

    if (newMode === FORWARD_MODE) {
      emailUseCase.toSelector.select([]);
      emailUseCase.ccSelector.select([]);
      emailUseCase.bccSelector.select([]);
      emailUseCase.updateEmailContent(`${event.content}`);
    }
    setMode(newMode);
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

    const isFormEmpty =
      !content.length || content === `<p class="my-3"><br></p>`;
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

  useKey('Escape', handleClosePreview);

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
          modal
          replyMode={mode}
          replyToId={event.id}
          emailUseCase={emailUseCase}
          onModeChange={handleModeChange}
          onDiscard={handleExitEditorAndCleanData}
        />

        <ConfirmDeleteDialog
          isOpen={isOpen}
          isLoading={false}
          colorScheme='primary'
          confirmButtonLabel='Send'
          label={`Send this email?`}
          emailUseCase={emailUseCase}
          cancelButtonLabel='Discard'
          onClose={handleExitEditorAndCleanData}
          onConfirm={() => emailUseCase.createEmail(event.id)}
          description={`You have typed an unsent email. Do you want to send it, or discard it?`}
        />
      </div>
    </TimelinePreviewBackdrop>
  );
};
