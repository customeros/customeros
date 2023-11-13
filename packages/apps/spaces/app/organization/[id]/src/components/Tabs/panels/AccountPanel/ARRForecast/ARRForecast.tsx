'use client';
import { useIsMutating, useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { UseDisclosureReturn } from '@ui/utils';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody } from '@ui/presentation/Card';
import { RenewalLikelihoodProbability } from '@graphql/types';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { RenewalForecastType } from '@organization/src/components/Tabs/panels/AccountPanel/RenewalForecast';
import { useUpdateRenewalLikelihoodMutation } from '@organization/src/graphql/updateRenewalLikelyhood.generated';

interface ARRForecastProps {
  name: string;
  isInitialLoading?: boolean;
  infoModal: UseDisclosureReturn;
  aRRForecast?: RenewalForecastType;
  renewalProbability?: RenewalLikelihoodProbability | null;
}

export const ARRForecast = ({
  isInitialLoading,
  aRRForecast,
  renewalProbability,
  infoModal,
  name,
}: ARRForecastProps) => {
  const isRestoring = useIsRestoring();

  const isMutating = useIsMutating({
    mutationKey: useUpdateRenewalLikelihoodMutation.getKey(),
  });

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        variant='ghost'
        bg='transparent'
        cursor='default'
        boxShadow='none'
        sx={{
          '& button': {
            opacity: 0,
            transition: 'opacity 0.2s linear',
          },
        }}
        _hover={{
          '& button': {
            opacity: 1,
          },
        }}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={
              aRRForecast?.amount && !aRRForecast?.updatedBy
                ? getFeatureIconColor(renewalProbability)
                : 'gray'
            }
          >
            <CurrencyDollar />
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
                  ARR forecast
                </Heading>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Help'
                  onClick={(e) => {
                    e.stopPropagation();
                    infoModal.onOpen();
                  }}
                  icon={<Icons.HelpCircle color='gray.400' />}
                />
              </Flex>
            </Flex>

            <Heading fontSize='2xl' color='gray.700'>
              {isMutating && (!isInitialLoading || !isRestoring)
                ? 'Calculating...'
                : formatCurrency(aRRForecast?.amount ?? 0)}
            </Heading>
          </Flex>
        </CardBody>
      </Card>
      <InfoDialog
        isOpen={infoModal.isOpen}
        onClose={infoModal.onClose}
        onConfirm={infoModal.onClose}
        confirmButtonLabel='Got it'
        label='ARR forecast'
      >
        <Text fontSize='sm' fontWeight='normal' mb={4}>
          Annual Recurring Revenue (ARR) is the total amount of money you can
          expect to receive from{' '}
          <Text as='span' fontWeight='medium'>
            {name}
          </Text>{' '}
          for the next 12 months.
        </Text>
        <Text fontSize='sm' fontWeight='normal'>
          It includes all renewals but excludes one-time services. Renewals are
          discounted based on the renewal likelihood.
        </Text>
      </InfoDialog>
    </>
  );
};
