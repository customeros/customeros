import { useParams } from 'next/navigation';
import React, { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { toastError } from '@ui/presentation/Toast';
import { Card, CardHeader } from '@ui/presentation/Card';
import { getDifferenceFromNow } from '@shared/util/date';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { useUpdateOpportunityRenewalMutation } from '@organization/src/graphql/updateOpportunityRenewal.generated';
import {
  Opportunity,
  InternalStage,
  ContractRenewalCycle,
  OpportunityRenewalLikelihood,
} from '@graphql/types';
import { RenewalDetailsModal } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/RenewalARR/RenewalDetailsModal';
import { useUpdateRenewalDetailsContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';

import {
  getRenewalLikelihoodColor,
  getRenewalLikelihoodLabel,
} from '../../utils';

interface RenewalARRCardProps {
  hasEnded: boolean;
  startedAt: string;
  opportunity: Opportunity;
  renewCycle: ContractRenewalCycle;
}
export const RenewalARRCard = ({
  startedAt,
  hasEnded,
  renewCycle,
  opportunity,
}: RenewalARRCardProps) => {
  const orgId = useParams()?.id as string;
  const queryClient = useQueryClient();
  const client = getGraphQLClient();

  const { modal } = useUpdateRenewalDetailsContext();
  const [isLocalOpen, setIsLocalOpen] = useState(false);

  const getContractsQueryKey = useGetContractsQuery.getKey({
    id: orgId,
  });
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const updateOpportunityMutation = useUpdateOpportunityRenewalMutation(
    client,
    {
      onMutate: ({ input }) => {
        queryClient.cancelQueries({ queryKey: getContractsQueryKey });

        queryClient.setQueryData<GetContractsQuery>(
          getContractsQueryKey,
          (currentCache) => {
            if (!currentCache || !currentCache?.organization) return;

            return produce(currentCache, (draft) => {
              if (draft?.['organization']?.['contracts']) {
                draft['organization']['contracts']?.map(
                  (contractData, index) => {
                    return (contractData.opportunities ?? []).map(
                      (opportunity) => {
                        const { opportunityId, ...rest } = input;
                        if ((opportunity as Opportunity).id === opportunityId) {
                          return {
                            ...opportunity,
                            ...rest,
                            renewalUpdatedByUserAt: new Date().toISOString(),
                          };
                        }

                        return opportunity;
                      },
                    );
                  },
                );
              }
            });
          },
        );
        const previousEntries =
          queryClient.getQueryData<GetContractsQuery>(getContractsQueryKey);

        return { previousEntries };
      },
      onError: (_, __, context) => {
        queryClient.setQueryData<GetContractsQuery>(
          getContractsQueryKey,
          context?.previousEntries,
        );
        toastError(
          'Failed to update renewal details',
          'update-renewal-details-error',
        );
      },
      onSettled: () => {
        modal.onClose?.();

        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }
        timeoutRef.current = setTimeout(() => {
          queryClient.invalidateQueries({ queryKey: getContractsQueryKey });
        }, 900);
      },
    },
  );

  const differenceInMonths = DateTimeUtils.differenceInMonths(
    new Date().toISOString(),
    startedAt,
  );
  const hasRenewed = startedAt
    ? renewCycle === ContractRenewalCycle.AnnualRenewal
      ? differenceInMonths > 12
      : differenceInMonths > 1
    : null;

  const formattedMaxAmount = formatCurrency(opportunity.maxAmount ?? 0);
  const formattedAmount = formatCurrency(hasEnded ? 0 : opportunity.amount);

  const hasRewenewChanged = formattedMaxAmount !== formattedAmount; // should be also less

  const hasRenewalLikelihoodZero =
    opportunity?.renewalLikelihood === OpportunityRenewalLikelihood.ZeroRenewal;
  const timeToRenewal = getDifferenceFromNow(opportunity.renewedAt).join(' ');

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <>
      <Card
        px='4'
        py='3'
        w='full'
        my={2}
        size='lg'
        variant='outline'
        cursor={hasEnded ? 'default' : 'pointer'}
        border='1px solid'
        borderColor='gray.200'
        position='relative'
        onClick={() => {
          if (opportunity?.internalStage === InternalStage.ClosedLost) return;
          modal.onOpen();
          setIsLocalOpen(true);
        }}
        width={hasRenewed ? 'calc(100% - .5rem)' : 'auto'}
        sx={
          hasRenewed
            ? {
                right: -2,
                '&:after': {
                  content: "''",
                  width: 2,
                  height: '80%',
                  left: '-9px',
                  top: '6px',
                  bg: 'white',
                  position: 'absolute',
                  borderTopLeftRadius: 'md',
                  borderBottomLeftRadius: 'md',
                  border: '1px solid',
                  borderColor: 'gray.200',
                },
              }
            : {}
        }
      >
        <CardHeader
          as={Flex}
          p='0'
          w='full'
          alignItems='center'
          justifyContent='center'
          gap={4}
        >
          <FeaturedIcon size='md' minW='10' colorScheme='primary'>
            <ClockFastForward />
          </FeaturedIcon>
          <Flex alignItems='center' justifyContent='space-between' w='full'>
            <Flex flexDir='column' gap='1px'>
              <Flex flex={1} alignItems='center'>
                <Heading size='sm' color='gray.700' noOfLines={1}>
                  Renewal ARR
                </Heading>

                {!hasEnded && opportunity.renewedAt && startedAt && (
                  <Text color='gray.500' ml={1} fontSize='sm'>
                    {timeToRenewal}
                  </Text>
                )}
              </Flex>

              {opportunity?.renewalLikelihood && (
                <Text w='full' color='gray.500' fontSize='sm' lineHeight={1}>
                  {!hasEnded ? (
                    <>
                      Likelihood{' '}
                      <Text
                        as='span'
                        fontWeight='medium'
                        color={`${getRenewalLikelihoodColor(
                          opportunity.renewalLikelihood as OpportunityRenewalLikelihood,
                        )}.500`}
                        textTransform='capitalize'
                      >
                        {getRenewalLikelihoodLabel(
                          opportunity?.renewalLikelihood as OpportunityRenewalLikelihood,
                        )}
                      </Text>
                    </>
                  ) : (
                    'Closed lost'
                  )}
                </Text>
              )}
            </Flex>

            <Flex flexDir='column'>
              <Text fontWeight='semibold'>{formattedAmount}</Text>

              {hasRewenewChanged && (
                <Text
                  fontSize='sm'
                  textAlign='right'
                  textDecoration='line-through'
                >
                  {formattedMaxAmount}
                </Text>
              )}
            </Flex>
          </Flex>
        </CardHeader>
      </Card>

      {hasRenewalLikelihoodZero ? (
        <InfoDialog
          isOpen={modal.isOpen && isLocalOpen}
          onClose={modal.onClose}
          onConfirm={modal.onClose}
          confirmButtonLabel='Got it'
          label='This contract ends soon'
        >
          <Text fontSize='sm' fontWeight='normal' mt={1}>
            The renewal likelihood has been downgraded to Zero because the
            contract is set to end within the current renewal cycle.
          </Text>
        </InfoDialog>
      ) : (
        <RenewalDetailsModal
          updateOpportunityMutation={updateOpportunityMutation}
          isOpen={modal.isOpen && isLocalOpen}
          onClose={() => {
            modal.onClose();
            setIsLocalOpen(false);
          }}
          data={opportunity}
        />
      )}
    </>
  );
};
