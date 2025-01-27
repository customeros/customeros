import { useParams } from 'react-router-dom';
import { VirtuosoHandle } from 'react-virtuoso';
import { useMemo, useState, useEffect } from 'react';

import { useKey } from 'rooks';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { SelectOption } from '@ui/utils/types.ts';
import { useStore } from '@shared/hooks/useStore';
import { useChannel } from '@shared/hooks/useChannel';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { EmailParticipant, InteractionEvent } from '@graphql/types';
import { getEmailParticipantsNameAndEmailSelection } from '@utils/getParticipantsName';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';
import { getEmailParticipantsByType } from '@organization/components/Timeline/PastZone/events/email/utils';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/PastZone/hooks/updateCacheWithNewEvent';
import { TimelinePreviewBackdrop } from '@organization/components/Timeline/shared/TimelineEventPreview/TimelinePreviewBackdrop';
import { ComposeEmailContainer } from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmailContainer';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import { ModeChangeButtons } from '@organization/components/Timeline/PastZone/events/email/compose-email/EmailResponseModeChangeButtons.tsx';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { EmailMetaDataEntry } from './EmailMetaDataEntry';

import postStamp from '/backgrounds/organization/post-stamp.webp';

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';

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
  const store = useStore();
  const [mode, setMode] = useState(REPLY_MODE);
  const { to, cc, bcc } = getEmailParticipantsByType(event?.sentTo || []);

  const fromParticipants = event?.sentBy || [];
  const from = getEmailParticipantsNameAndEmailSelection(fromParticipants);

  const handleEmailSendSuccess = (response: unknown) => {
    updateTimelineCache(response, ['GetTimeline.infinite']);
    invalidateQuery();
  };
  const { currentUserId } = useChannel(`finder:${store.session.value.tenant}`);

  const attendees = (event?.interactionSession?.attendedBy
    ?.map((i) => i as EmailParticipant)
    .map((i) => i.emailParticipant.email) ?? []) as string[];
  const emailUseCase = useMemo(
    () =>
      new TimelineEmailUsecase(
        orgId,
        attendees,
        currentUserId ?? '',
        handleEmailSendSuccess,
      ),
    [orgId, event.id, isOpen, currentUserId],
  );

  useEffect(() => {
    if (emailUseCase) {
      emailUseCase.updateSubject(`Re: ${subject}`);

      if (from?.[0].value) {
        emailUseCase.fromSelector.select(from[0] as SelectOption);
      }
      emailUseCase.toSelector.select(
        getEmailParticipantsNameAndEmailSelection(to) || [],
      );
      emailUseCase.ccSelector.select(
        getEmailParticipantsNameAndEmailSelection(cc) || [],
      );
      emailUseCase.bccSelector.select(
        getEmailParticipantsNameAndEmailSelection(bcc) || [],
      );
    }
  }, [emailUseCase]);

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
          ...getEmailParticipantsNameAndEmailSelection([
            ...to.filter(
              (e) =>
                e.emailParticipant.email !== store.session.value.profile.email,
            ),
          ]),
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
        ...getEmailParticipantsNameAndEmailSelection([
          ...cc,
          ...to.filter(
            (e) =>
              e.emailParticipant.email !== store.session.value.profile.email,
          ),
        ]),
      ];
      const newBCC = getEmailParticipantsNameAndEmailSelection(bcc);

      emailUseCase.toSelector.select(newTo as SelectOption[]);
      emailUseCase.ccSelector.select(
        removeDuplicates(newTo, newCC) as SelectOption[],
      );
      emailUseCase.bccSelector.select(newBCC as SelectOption[]);
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
    emailUseCase.resetEditor();
    onClose();
    closeModal();
  };

  const handleClosePreview = (): void => {
    const showConfirmationDialog = emailUseCase.canExitSafely;

    if (showConfirmationDialog) {
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
              <img
                alt='Email'
                src={postStamp}
                className='min-w-[48px] max-w-[48px] h-[70px]'
              />
            </div>
          </div>

          {event?.content && (
            <HtmlContentRenderer htmlContent={event.content} />
          )}
        </div>

        <div className='rounded-b-2xl py-4 px-6 overflow-visible pt-1 border-dashed border-t-[1px] border-gray-200 bg-grayBlue-50 rounded-none max-h-[50vh]'>
          <div style={{ position: 'relative' }}>
            <ModeChangeButtons handleModeChange={handleModeChange} />
          </div>
          <ComposeEmailContainer
            modal
            key={mode}
            replyToId={event.id}
            emailUseCase={emailUseCase}
            onClose={handleClosePreview}
            onDiscard={handleExitEditorAndCleanData}
          />
        </div>

        <ConfirmDeleteDialog
          isOpen={isOpen}
          isLoading={false}
          colorScheme='primary'
          confirmButtonLabel='Send'
          label={`Send this email?`}
          cancelButtonLabel='Discard'
          onClose={handleExitEditorAndCleanData}
          onConfirm={() => emailUseCase.createEmail(event.id)}
          description={`You have typed an unsent email. Do you want to send it, or discard it?`}
        />
      </div>
    </TimelinePreviewBackdrop>
  );
};
