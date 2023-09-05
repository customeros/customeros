import React from 'react';
import { Action } from '@graphql/types';
import { Card, CardFooter } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { TimelineEventPreviewHeader } from '../../../preview/header/TimelineEventPreviewHeader';
import { useTimelineEventPreviewContext } from '../../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { CardBody } from '@chakra-ui/card';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { getFeatureIconColor } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { Divider } from '@ui/presentation/Divider';
import {
  getCurrencyString,
  getMetadata,
} from '@organization/components/Timeline/events/action/utils';

export const RenewalForecastUpdatedActionPreview = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const event = modalContent as Action;
  const metadata = getMetadata(event?.metadata);
  const forecastedAmount = event.content && getCurrencyString(event.content);
  const [_, author] = event.content?.split('by ') ?? [];
  const isCreatedBySystem = event.content?.includes('default');
  const colorScheme =
    forecastedAmount && isCreatedBySystem
      ? getFeatureIconColor(metadata.likelihood)
      : 'gray';

  const getForecastMetaInfo = () => {
    if (isCreatedBySystem) {
      return 'Calculated from billing amount';
    }

    return `Set by ${author} ${DateTimeUtils.timeAgo(modalContent?.createdAt, {
      addSuffix: true,
    })}`;
  };
  return (
    <>
      <TimelineEventPreviewHeader
        date={modalContent?.createdAt}
        name={'Renewal forecast'}
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <CardBody
        as={Flex}
        justify='space-between'
        px='6'
        py='0'
        paddingBottom='6'
        overflowY='auto'
        maxH='calc(100vh - 4rem - 56px - 51px - 16px - 16px);'
      >
        <Card mt={3} p='4' w='full' size='lg' boxShadow='xs' variant='outline'>
          <CardBody as={Flex} p='0' align='center'>
            <FeaturedIcon size='md' minW='10' colorScheme={colorScheme}>
              <Icons.Calculator />
            </FeaturedIcon>
            <Flex
              ml='5'
              w='full'
              align='center'
              columnGap={4}
              justify='space-between'
            >
              <Flex flexDir='column'>
                <Flex align='center'>
                  <Heading
                    size='sm'
                    whiteSpace='nowrap'
                    fontWeight='semibold'
                    color='gray.700'
                    mr={2}
                  >
                    Renewal forecast
                  </Heading>
                </Flex>
                <Text fontSize='xs' color='gray.500'>
                  {getForecastMetaInfo()}
                </Text>
              </Flex>

              <Heading fontSize='2xl' color='gray.700'>
                {forecastedAmount}
              </Heading>
            </Flex>
          </CardBody>
          {!isCreatedBySystem && forecastedAmount && (
            <CardFooter p='0' as={Flex} flexDir='column'>
              <Divider mt='4' mb='2' />
              <Flex align='flex-start'>
                {metadata?.reason ? (
                  <Icons.File2 color='gray.400' />
                ) : (
                  <Icons.FileCross viewBox='0 0 16 16' color='gray.400' />
                )}

                <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                  {metadata?.reason || 'No reason provided'}
                </Text>
              </Flex>
            </CardFooter>
          )}
        </Card>
      </CardBody>
    </>
  );
};
