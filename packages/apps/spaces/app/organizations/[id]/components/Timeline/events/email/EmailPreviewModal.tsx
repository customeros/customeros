import React, { useRef } from 'react';
import { Card, CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import { IconButton } from '@ui/form/IconButton';
import Image from 'next/image';
import styles from './EmailPreviewModal.module.scss';
import { EmailMetaDataEntry } from './EmailMetaDataEntry';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { useCopyToClipboard } from '@spaces/hooks/useCopyToClipboard';
import sanitizeHtml from 'sanitize-html';
import { DateTimeUtils } from '@spaces/utils/date';
import { getEmailParticipantsByType } from '@organization/components/Timeline/events/email/utils';
import CopyLink from '@spaces/atoms/icons/CopyLink';
import Times from '@spaces/atoms/icons/Times';
import { useOutsideClick } from '@spaces/hooks/useOutsideClick';
import { ComposeEmail } from '@organization/components/Timeline/events/email/compose-email/ComposeEmail';
import { getEmailParticipantsNameAndEmail } from '@spaces/utils/getParticipantsName';

export const EmailPreviewModal: React.FC = () => {
  const { closeModal, isModalOpen, modalContent } =
    useTimelineEventPreviewContext();
  const ref = useRef(null);
  const [_, copy] = useCopyToClipboard();
  useOutsideClick({
    ref: ref,
    handler: () => closeModal(),
  });
  if (!isModalOpen || !modalContent) {
    return null;
  }
  const { to, cc, bcc } = getEmailParticipantsByType(modalContent.sentTo);

  return (
    <div className={styles.backdrop}>
      <ScaleFade initialScale={0.9} in={isModalOpen} unmountOnExit>
        <Card
          ref={ref}
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
                <Tooltip
                  label='Copy link'
                  aria-label='Share'
                  placement='bottom'
                >
                  <IconButton
                    variant='ghost'
                    aria-label='Close preview'
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
                    onClick={() => closeModal()}
                  />
                </Tooltip>
              </Flex>
            </Flex>
          </CardHeader>

          <CardBody mt={0} maxHeight='50%' overflow='auto' pb={6}>
            <Flex direction='row' justify='space-between' mb={3}>
              <Flex direction='column' align='flex-start'>
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
                  src={'/backgrounds/organization/poststamp1.webp'}
                  alt='Email'
                  width={54}
                  height={70}
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
            to={getEmailParticipantsNameAndEmail(to, 'value')}
            cc={getEmailParticipantsNameAndEmail(cc, 'value')}
            bcc={getEmailParticipantsNameAndEmail(bcc, 'value')}
            from={getEmailParticipantsNameAndEmail(
              modalContent.sentBy,
              'value',
            )}
            subject={modalContent?.interactionSession?.name || ''}
            emailContent={modalContent.content || ''}
          />
        </Card>
      </ScaleFade>
    </div>
  );
};
