import React from 'react';
import { Action } from '@graphql/types';
import { Card, CardFooter } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { TimelineEventPreviewHeader } from '../../../preview/header/TimelineEventPreviewHeader';
import { CardBody } from '@chakra-ui/card';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { Divider } from '@ui/presentation/Divider';
import {
  getCurrencyString,
  getMetadata,
} from '@organization/src/components/Timeline/events/action/utils';
import {
  useTimelineEventPreviewMethodsContext,
  useTimelineEventPreviewStateContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

export const RenewalForecastUpdatedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as Action;
  const metadata = getMetadata(event?.metadata);
  const forecastedAmount = event.content && getCurrencyString(event.content);
  const isCreatedBySystem = event.content?.includes('default');
  const colorScheme =
    forecastedAmount && isCreatedBySystem
      ? getFeatureIconColor(metadata.likelihood)
      : 'gray';

  const [preText, postText] = event.content?.split('by ') ?? [];
  const authorText = isCreatedBySystem ? event.content : `${preText} by`;

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
            <Text
              my={1}
              maxW='500px'
              noOfLines={2}
              ml={2}
              fontSize='sm'
              color='gray.700'
            >
              {authorText}
              {!isCreatedBySystem && (
                <Text color='gray.500' as='span' ml={1}>
                  {postText}
                </Text>
              )}
            </Text>
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
