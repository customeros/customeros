'use client';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { useDisclosure } from '@ui/utils';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';

import {
  RenewalForecastModal,
  Value as RenewalForecastValue,
} from './RenewalForecastModal';
import { FC, useState } from 'react';

export const RenewalForecast: FC = () => {
  const [renewalForecast, setRenewalForecast] = useState<RenewalForecastValue>({
    reason: '',
    forecast: '',
  });
  const update = useDisclosure();
  const info = useDisclosure();
  const { forecast, reason } = renewalForecast;

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
          <FeaturedIcon size='md' colorScheme={forecast ? 'success' : 'gray'}>
            <Icons.Calculator />
          </FeaturedIcon>
          <Flex ml='5' align='center' justify='space-between' w='full'>
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading size='sm' fontWeight='semibold' color='gray.700'>
                  Renewal Forecast
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
                {!forecast ? 'Not calculated yet' : 'Set by Unknown just now'}
              </Text>
            </Flex>

            <Heading fontSize='2xl' color={!forecast ? 'gray.400' : 'gray.700'}>
              {!forecast
                ? 'Unknown'
                : Intl.NumberFormat('en-US', {
                    style: 'currency',
                    currency: 'USD',
                  }).format(parseFloat(forecast))}
            </Heading>
          </Flex>
        </CardBody>
        <CardFooter p='0' as={Flex} flexDir='column'>
          {reason && (
            <>
              <Divider mt='4' mb='2' />
              <Flex align='flex-start'>
                <Icons.File2 color='gray.400' />
                <Text color='gray.500' fontSize='xs' ml='1'>
                  {reason}
                </Text>
              </Flex>
            </>
          )}
        </CardFooter>
      </Card>

      <RenewalForecastModal
        value={renewalForecast}
        onChange={setRenewalForecast}
        isOpen={update.isOpen}
        onClose={update.onClose}
      />

      <InfoDialog
        isOpen={info.isOpen}
        onClose={info.onClose}
        onConfirm={info.onClose}
        confirmButtonLabel='Got it'
        label='Renewal likelihood'
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
