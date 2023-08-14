'use client';
import { FC } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { useDisclosure } from '@ui/utils';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';

import { RenewalForecastModal } from './RenewalForecastModal';
import { RenewalForecast as RenewalForecastT } from '@graphql/types';
import { getUserDisplayData } from '@spaces/utils/getUserEmail';
import { DateTimeUtils } from '@spaces/utils/date';

export type RenewalForecastType = RenewalForecastT & { amount?: string | null };

export const RenewalForecast: FC<{
  renewalForecast: RenewalForecastType;
  name: string;
}> = ({ renewalForecast, name }) => {
  const update = useDisclosure();
  const info = useDisclosure();
  const { amount, comment } = renewalForecast;

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        boxShadow='xs'
        variant='outline'
        cursor='pointer'
        onClick={update.onOpen}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon size='md' colorScheme={amount ? 'success' : 'gray'}>
            <Icons.Calculator />
          </FeaturedIcon>
          <Flex
            ml='5'
            align='center'
            justify='space-between'
            w='full'
            columnGap={4}
          >
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading size='sm' fontWeight='semibold' color='gray.700' mr={2}>
                  Renewal forecast
                </Heading>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Help'
                  onClick={(e) => {
                    e.stopPropagation();
                    info.onOpen();
                  }}
                  icon={<Icons.HelpCircle color='gray.400' />}
                />
              </Flex>
              <Text fontSize='xs' color='gray.500'>
                {!amount
                  ? 'Not calculated yet'
                  : `Set by 
                ${getUserDisplayData(renewalForecast?.updatedBy)}
                 ${DateTimeUtils.timeAgo(renewalForecast.updatedAt, {
                   addSuffix: true,
                 })}`}
              </Text>
            </Flex>

            <Heading fontSize='2xl' color={!amount ? 'gray.400' : 'gray.700'}>
              {!amount
                ? 'Unknown'
                : Intl.NumberFormat('en-US', {
                    style: 'currency',
                    currency: 'USD',
                    minimumFractionDigits: 0,
                  }).format(parseFloat(`${amount}`))}
            </Heading>
          </Flex>
        </CardBody>
        <CardFooter p='0' as={Flex} flexDir='column'>
          {comment && (
            <>
              <Divider mt='4' mb='2' />
              <Flex align='flex-start'>
                <Icons.File2 color='gray.400' />
                <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                  {comment}
                </Text>
              </Flex>
            </>
          )}
        </CardFooter>
      </Card>

      <RenewalForecastModal
        renewalForecast={{
          amount: renewalForecast.amount,
          comment: renewalForecast.comment,
        }}
        name={name}
        isOpen={update.isOpen}
        onClose={update.onClose}
      />

      <InfoDialog
        isOpen={info.isOpen}
        onClose={info.onClose}
        onConfirm={info.onClose}
        confirmButtonLabel='Got it'
        label='Renewal forecast'
      >
        <Text fontSize='sm' fontWeight='normal'>
          The renewal forecast gives you a way to roughly project revenue per
          customer and across your entire portfolio.
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          {`It's calculated by discounting the subscription amount based on the
          renewal likelihood—Medium, Low, or Zero.`}
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          You can override this forecast at any time.
        </Text>
      </InfoDialog>
    </>
  );
};
