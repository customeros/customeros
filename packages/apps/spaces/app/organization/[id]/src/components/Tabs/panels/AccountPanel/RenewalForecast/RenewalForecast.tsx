'use client';
import { useIsMutating, useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { UseDisclosureReturn } from '@ui/utils';
import { Heading } from '@ui/typography/Heading';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { DateTimeUtils } from '@spaces/utils/date';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { getUserDisplayData } from '@spaces/utils/getUserEmail';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import {
  RenewalLikelihoodProbability,
  RenewalForecast as RenewalForecastT,
} from '@graphql/types';
import { useUpdateRenewalLikelihoodMutation } from '@organization/src/graphql/updateRenewalLikelyhood.generated';

import { RenewalForecastModal } from './RenewalForecastModal';

export type RenewalForecastType = RenewalForecastT & { amount?: string | null };

interface RenewalForecastProps {
  name: string;
  isInitialLoading?: boolean;
  infoModal: UseDisclosureReturn;
  updateModal: UseDisclosureReturn;
  renewalForecast: RenewalForecastType;
  renewalProbability?: RenewalLikelihoodProbability | null;
}

export const RenewalForecast = ({
  isInitialLoading,
  renewalForecast,
  renewalProbability,
  name,
  infoModal,
  updateModal,
}: RenewalForecastProps) => {
  const isRestoring = useIsRestoring();

  const isMutating = useIsMutating({
    mutationKey: useUpdateRenewalLikelihoodMutation.getKey(),
  });

  const getForecastMetaInfo = () => {
    if (!renewalForecast?.amount) {
      return 'Not calculated yet';
    }

    if (!renewalForecast?.updatedBy) {
      return 'Calculated from billing amount';
    }

    return `Set by ${getUserDisplayData(
      renewalForecast?.updatedBy,
    )} ${DateTimeUtils.timeAgo(renewalForecast?.updatedAt, {
      addSuffix: true,
    })}`;
  };

  const isAmountSet =
    renewalForecast?.amount !== null && renewalForecast?.amount !== undefined;

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        variant='outline'
        cursor='pointer'
        boxShadow='xs'
        _hover={{
          boxShadow: 'md',
        }}
        transition='all 0.2s ease-out'
        onClick={updateModal.onOpen}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={
              renewalForecast?.amount && !renewalForecast?.updatedBy
                ? getFeatureIconColor(renewalProbability)
                : 'gray'
            }
          >
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
              <Text fontSize='xs' color='gray.500'>
                {getForecastMetaInfo()}
              </Text>
            </Flex>

            <Heading
              fontSize='2xl'
              color={isAmountSet ? 'gray.700' : 'gray.400'}
            >
              {isMutating && (!isInitialLoading || !isRestoring)
                ? 'Calculating...'
                : isAmountSet
                ? formatCurrency(renewalForecast?.amount ?? 0)
                : 'Unknown'}
            </Heading>
          </Flex>
        </CardBody>
        {!!renewalForecast?.amount && renewalForecast?.updatedBy && (
          <CardFooter p='0' as={Flex} flexDir='column'>
            <Divider mt='4' mb='2' />
            <Flex align='flex-start'>
              {renewalForecast?.comment ? (
                <Icons.File2 color='gray.400' />
              ) : (
                <Icons.FileCross viewBox='0 0 16 16' color='gray.400' />
              )}

              <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                {renewalForecast?.comment || 'No reason provided'}
              </Text>
            </Flex>
          </CardFooter>
        )}
      </Card>

      <RenewalForecastModal
        renewalForecast={{
          amount: renewalForecast?.amount,
          comment: renewalForecast?.comment,
        }}
        renewalProbability={renewalProbability}
        name={name}
        isOpen={updateModal.isOpen}
        onClose={updateModal.onClose}
      />

      <InfoDialog
        isOpen={infoModal.isOpen}
        onClose={infoModal.onClose}
        onConfirm={infoModal.onClose}
        confirmButtonLabel='Got it'
        label='ARR forecast'
      >
        <Text fontSize='sm' fontWeight='normal'>
          The ARR forecast gives you a way to roughly project revenue per
          customer and across your entire portfolio.
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          {`It's calculated by discounting the renewal potential (billing amount * billings per cycle) based on the renewal likelihood—Medium, Low, or Zero.`}
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          You can override this forecast at any time.
        </Text>
      </InfoDialog>
    </>
  );
};
