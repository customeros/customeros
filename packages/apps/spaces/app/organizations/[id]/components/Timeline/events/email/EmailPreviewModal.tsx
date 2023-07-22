import React from 'react';
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
import TimesOutlined from '@spaces/atoms/icons/TimesOutlined';
import Copy from '@spaces/atoms/icons/Copy';
import { useCopyToClipboard } from '@spaces/hooks/useCopyToClipboard';
import sanitizeHtml from 'sanitize-html';
import { InteractionEventParticipant } from '@graphql/types';
import { DateTimeUtils } from '@spaces/utils/date';

export const EmailPreviewModal: React.FC = () => {
  const { closeModal, isModalOpen, modalContent } =
    useTimelineEventPreviewContext();
  const [_, copy] = useCopyToClipboard();

  if (!isModalOpen || !modalContent) {
    return null;
  }
  const cc = (modalContent.sentTo || []).filter(
    (e: InteractionEventParticipant) => e.type === 'CC',
  );
  const bcc = (modalContent.sentTo || []).filter(
    (e: InteractionEventParticipant) => e.type === 'BCC',
  );
  const to = (modalContent.sentTo || []).filter(
    (e: InteractionEventParticipant) => e.type === 'TO',
  );
  return (
    <div className={styles.backdrop}>
      <ScaleFade initialScale={0.9} in={isModalOpen} unmountOnExit>
        <Card
          borderRadius='xl'
          height='100%'
          maxHeight='calc(100vh - 6rem)'
          overflow='scroll'
        >
          <CardHeader pb={1} position='sticky' background='white' top={0}>
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
                  aria-label='Copy link'
                  placement='top'
                >
                  <IconButton
                    variant='ghost'
                    aria-label='Close preview'
                    color='gray.500'
                    borderRadius={30}
                    padding={0}
                    height='24px'
                    width='24px'
                    minWidth='24px'
                    mr={3}
                    icon={<Copy color='#98A2B3' height='24px' />}
                    onClick={() => copy(window.location.href)}
                  />
                </Tooltip>
                <Tooltip label='Close' aria-label='close' placement='top'>
                  <IconButton
                    variant='ghost'
                    aria-label='Close preview'
                    color='gray.500'
                    borderRadius={30}
                    padding={0}
                    height='24px'
                    width='24px'
                    minWidth='24px'
                    icon={<TimesOutlined color='#98A2B3' height='24px' />}
                    onClick={() => closeModal()}
                  />
                </Tooltip>
              </Flex>
            </Flex>
          </CardHeader>

          <CardBody mt={0}>
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
                  src={'/backgrounds/organization/poststamp.webp'}
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
        </Card>
      </ScaleFade>
    </div>
  );
};
