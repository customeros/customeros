import React, { useState } from 'react';
import { Card, CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import { IconButton } from '@ui/form/IconButton';
import styles from './EmailPreviewModal.module.scss';
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
import { convert } from 'html-to-text';
import { ConfirmDeleteDialog } from '@ui/presentation/Modal/ConfirmDeleteDialog';
import { useDisclosure } from '@chakra-ui/react-use-disclosure';

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';

interface EmailPreviewModalProps {
  invalidateQuery: () => void;
}
export const EmailPreviewModal: React.FC<EmailPreviewModalProps> = ({
  invalidateQuery,
}) => {
  const { closeModal, isModalOpen, modalContent } =
    useTimelineEventPreviewContext();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const subject = modalContent?.interactionSession?.name || '';
  const [_, copy] = useCopyToClipboard();
  const searchParams = useSearchParams();
  const { data: session } = useSession();
  const [mode, setMode] = useState(REPLY_MODE);
  const [isSending, setIsSending] = useState(false);
  const { to, cc, bcc } = getEmailParticipantsByType(
    modalContent?.sentTo || [],
  );
  const from = getEmailParticipantsNameAndEmail(
    modalContent?.sentBy || [],
    'value',
  );
  const formId = 'compose-email-preview-modal';

  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: getEmailParticipantsNameAndEmail(to, 'value'),
    cc: getEmailParticipantsNameAndEmail(cc, 'value'),
    bcc: getEmailParticipantsNameAndEmail(bcc, 'value'),
    subject: `Re: ${subject}`,
    content: '',
  });
  const { state, setDefaultValues, reset } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,

    stateReducer: (state, action, next) => {
      return next;
    },
  });

  if (!isModalOpen || !modalContent) {
    return null;
  }
  const handleEmailSendSuccess = () => {
    invalidateQuery()
    setIsSending(false);
    reset();
    closeModal();
  };
  const handleEmailSendError = () => {
    setIsSending(false);
  };
  const text = convert(modalContent?.content || '', {
    preserveNewlines: false,
    selectors: [
      {
        selector: 'a',
        options: { hideLinkHrefIfSameAsText: true, ignoreHref: true },
      },
    ],
  });

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
      newDefaultValues = new ComposeEmailDto({
        to: [...from, ...to],
        cc,
        bcc,
        subject: `Re: ${subject}`,
        content: mode === FORWARD_MODE ? '' : state.values.content,
      });
    }
    if (newMode === FORWARD_MODE) {
      newDefaultValues = new ComposeEmailDto({
        to: [],
        cc: [],
        bcc: [],
        subject: `Re: ${subject}`,
        content: `${state.values.content}\n ${text}`,
      });
    }
    setMode(newMode);
    setDefaultValues(newDefaultValues);
  };

  const handleClosePreview = () => {
    const isFormPristine = Object.values(state.fields)?.every(
      (e) => e.meta.pristine,
    );
    const isFormEmpty = Object.values(state.values)?.every((e) => !e.length);

    const showConfirmationDialog = !isFormPristine && !isFormEmpty;
    if (showConfirmationDialog) {
      onOpen();
    } else {
      closeModal();
    }
  };

  const handleSubmit = () => {
    const destination = [
      ...state.values.to,
      ...state.values.cc,
      ...state.values.bcc,
    ].map(({ value }) => value);
    const params = new URLSearchParams(searchParams ?? '');

    setIsSending(true);
    const id = params.get('events');
    return handleSendEmail(
      state.values.content,
      destination,
      id,
      state.values.subject,
      handleEmailSendSuccess,
      handleEmailSendError,
      session?.user?.email,
    );
  };
  return (
    <div className={styles.container}>
      <div
        className={styles.backdrop}
        onClick={() => (isModalOpen ? handleClosePreview() : null)}
      />
      <ScaleFade initialScale={0.9} in={isModalOpen} unmountOnExit>
        <Card
          zIndex={7}
          borderRadius='xl'
          height='100%'
          maxHeight='calc(100vh - 6rem)'
        >
          <CardHeader
            pb={1}
            position='sticky'
            background='white'
            top={0}
            borderRadius='xl'
          >
            <Flex
              direction='row'
              justifyContent='space-between'
              alignItems='center'
            >
              <div>
                <Heading size='sm' mb={2}>
                  {modalContent.interactionSession?.name}
                </Heading>
                <Text size='2xs' color='gray.500' fontSize='12px'>
                  {DateTimeUtils.format(
                    // @ts-expect-error this is correct (alias)
                    modalContent.date,
                    DateTimeUtils.dateWithHour,
                  )}
                </Text>
              </div>
              <Flex
                direction='row'
                justifyContent='flex-end'
                alignItems='center'
              >
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
                <EmailMetaDataEntry entryType='To' content={to} />
                {!!cc.length && (
                  <EmailMetaDataEntry entryType='CC' content={cc} />
                )}
                {!!bcc.length && (
                  <EmailMetaDataEntry entryType='BCC' content={bcc} />
                )}
                <EmailMetaDataEntry
                  entryType='Subject'
                  content={modalContent?.interactionSession?.name || ''}
                />
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
              {modalContent?.content && (
                <div
                  className={styles.normalize_email}
                  dangerouslySetInnerHTML={{
                    __html: sanitizeHtml(modalContent.content),
                  }}
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
          />
        </Card>
        <ConfirmDeleteDialog
          label='Discard this email?'
          description='Saving draft emails is not possible at the moment. Would you like to continue to discard this email?'
          confirmButtonLabel='Discard email'
          isOpen={isOpen}
          onClose={onClose}
          onConfirm={closeModal}
          isLoading={false}
        />
      </ScaleFade>
    </div>
  );
};
