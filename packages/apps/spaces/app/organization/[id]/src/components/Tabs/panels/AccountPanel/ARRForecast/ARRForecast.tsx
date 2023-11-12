'use client';
import { useIsMutating, useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { Card, CardBody } from '@ui/presentation/Card';
import { RenewalLikelihoodProbability } from '@graphql/types';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { RenewalForecastType } from '@organization/src/components/Tabs/panels/AccountPanel/RenewalForecast';
import { useUpdateRenewalLikelihoodMutation } from '@organization/src/graphql/updateRenewalLikelyhood.generated';

interface ARRForecastProps {
  name: string;
  isInitialLoading?: boolean;
  aRRForecast?: RenewalForecastType;
  renewalProbability?: RenewalLikelihoodProbability | null;
}

export const ARRForecast = ({
  isInitialLoading,
  aRRForecast,
  renewalProbability,
}: ARRForecastProps) => {
  const isRestoring = useIsRestoring();

  const isMutating = useIsMutating({
    mutationKey: useUpdateRenewalLikelihoodMutation.getKey(),
  });

  const isAmountSet =
    aRRForecast?.amount !== null && aRRForecast?.amount !== undefined;

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
              </Flex>
            </Flex>

            <Heading
              fontSize='2xl'
              color={isAmountSet ? 'gray.700' : 'gray.400'}
            >
              {isMutating && (!isInitialLoading || !isRestoring)
                ? 'Calculating...'
                : isAmountSet
                ? formatCurrency(aRRForecast?.amount ?? 0)
                : 'Unknown'}
            </Heading>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
