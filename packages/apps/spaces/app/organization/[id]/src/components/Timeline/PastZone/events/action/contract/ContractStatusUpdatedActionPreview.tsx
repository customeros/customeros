import React, { useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { getMetadata } from '@organization/src/components/Timeline/PastZone/events/action/utils';
import { iconsByStatus } from '@organization/src/components/Timeline/PastZone/events/action/contract/utils';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const ContractStatusUpdatedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as Action;
  const metadata = useMemo(() => {
    return getMetadata(event?.metadata);
  }, [event?.metadata]);
  const status = metadata?.status?.toLowerCase();
  // todo remove when contract name is passed from BE in metadata
  const getContractName = () => {
    const content = event.content ?? '';
    const endIndex =
      content.lastIndexOf(' is now live') > -1
        ? content.lastIndexOf(' is now live')
        : content.lastIndexOf(' renewed') > -1
        ? content.lastIndexOf(' renewed')
        : content.lastIndexOf(' has ended') > -1
        ? content.lastIndexOf(' has ended')
        : content.length;

    return content.substring(0, endIndex).trim();
  };
  const content = (event.content ?? '').substring(
    0,
    event?.content?.lastIndexOf(' '),
  );

  return (
    <>
      <TimelineEventPreviewHeader
        date={event?.createdAt}
        name={`${getContractName()} ${iconsByStatus[status].text} ${status}`}
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <Card m={6} mt={3} p='4' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={iconsByStatus[status].colorScheme as string}
          >
            {iconsByStatus[status].icon}
          </FeaturedIcon>
          <Text
            my={1}
            maxW='500px'
            noOfLines={2}
            ml={2}
            fontSize='sm'
            color='gray.700'
          >
            {content}
            <Text
              as='span'
              fontWeight={status === 'renewed' ? 'normal' : 'semibold'}
              ml={1}
            >
              {status}
            </Text>
          </Text>
        </CardBody>
        {metadata?.comment && (
          <CardFooter
            as={Flex}
            p='0'
            pt='3'
            mt='4'
            align='center'
            borderTop='1px solid'
            borderTopColor='gray.200'
          >
            <File02 color='gray.400' />
            <Text
              maxW='500px'
              noOfLines={2}
              ml={2}
              fontSize='sm'
              color='gray.500'
            >
              {metadata.content}
            </Text>
          </CardFooter>
        )}
      </Card>
    </>
  );
};
