import { useParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@spaces/utils/date';
import { toastError } from '@ui/presentation/Toast';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { Card, CardHeader } from '@ui/presentation/Card/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog/InfoDialog';
import {
  Opportunity,
  InternalStage,
  OpportunityRenewalLikelihood,
} from '@graphql/types';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import { useUpdateOpportunityRenewalMutation } from '@organization/graphql/updateOpportunityRenewal.generated';
import { useUpdateRenewalDetailsContext } from '@organization/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import { RenewalDetailsModal } from '@organization/components/Tabs/panels/AccountPanel/ContractOld/RenewalARR/RenewalDetailsModal';

import { getRenewalLikelihoodLabel } from '../../utils';

interface RenewalARRCardProps {
  hasEnded: boolean;
  startedAt: string;
  opportunity: Opportunity;
  currency?: string | null;
}
export const RenewalARRCard = ({
  startedAt,
  hasEnded,
  opportunity,
  currency,
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
                draft['organization']['contracts']?.map((contractData) => {
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
                });
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

  const formattedMaxAmount = formatCurrency(
    opportunity.maxAmount ?? 0,
    2,
    currency || 'USD',
  );
  const formattedAmount = formatCurrency(
    hasEnded ? 0 : opportunity.amount,
    2,
    currency || 'USD',
  );

  const hasRewenewChanged = formattedMaxAmount !== formattedAmount; // should be also less

  const hasRenewalLikelihoodZero =
    opportunity?.renewalLikelihood === OpportunityRenewalLikelihood.ZeroRenewal;
  const timeToRenewal = DateTimeUtils.getDifferenceFromNow(
    opportunity.renewedAt,
  ).join(' ');

  const showTimeToRenewal =
    !hasEnded &&
    opportunity.renewedAt &&
    startedAt &&
    !DateTimeUtils.isPast(opportunity.renewedAt);

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
        className={cn(
          'px-4 py-3 w-full my-2 border border-gray-200 relative bg-white rounded-lg shadow-xs',
          {
            'cursor-pointer': !hasEnded,
            'cursor-default': hasEnded,
          },
        )}
        onClick={() => {
          if (opportunity?.internalStage === InternalStage.ClosedLost) return;
          modal.onOpen();
          setIsLocalOpen(true);
        }}
      >
        <CardHeader className='flex items-center justify-between w-full gap-4'>
          <FeaturedIcon size='md' colorScheme='primary' className='ml-2 mr-2'>
            <ClockFastForward />
          </FeaturedIcon>
          <div className='flex items-center justify-between w-full'>
            <div className='flex flex-col'>
              <div className='flex flex-1 items-center'>
                <h1 className='text-gray-700 font-semibold text-sm line-height-1'>
                  Renewal ARR
                </h1>

                {showTimeToRenewal && (
                  <p className='ml-1 text-gray-500 text-sm inline'>
                    {timeToRenewal}
                  </p>
                )}
              </div>

              {opportunity?.renewalLikelihood && (
                <p className='w-full text-gray-500 text-sm line-height-1'>
                  {!hasEnded ? (
                    <>
                      Likelihood{' '}
                      <span
                        className={cn(`capitalize font-medium text-gray-500`, {
                          'text-greenLight-500':
                            opportunity?.renewalLikelihood ===
                            OpportunityRenewalLikelihood.HighRenewal,
                          'text-orangeDark-800':
                            opportunity?.renewalLikelihood ===
                            OpportunityRenewalLikelihood.LowRenewal,
                          'text-yellow-500':
                            opportunity?.renewalLikelihood ===
                            OpportunityRenewalLikelihood.MediumRenewal,
                        })}
                      >
                        {getRenewalLikelihoodLabel(
                          opportunity?.renewalLikelihood as OpportunityRenewalLikelihood,
                        )}
                      </span>
                    </>
                  ) : (
                    'Closed lost'
                  )}
                </p>
              )}
            </div>

            <div>
              <p className='font-semibold'>{formattedAmount}</p>

              {hasRewenewChanged && (
                <p className='text-sm text-right line-through'>
                  {formattedMaxAmount}
                </p>
              )}
            </div>
          </div>
        </CardHeader>
      </Card>

      {hasRenewalLikelihoodZero ? (
        <InfoDialog
          isOpen={modal.open && isLocalOpen}
          onClose={modal.onClose}
          onConfirm={modal.onClose}
          confirmButtonLabel='Got it'
          label='This contract ends soon'
          description=' The renewal likelihood has been downgraded to Zero because the
          contract is set to end within the current renewal cycle.'
        />
      ) : (
        <RenewalDetailsModal
          currency={currency}
          updateOpportunityMutation={updateOpportunityMutation}
          isOpen={modal.open && isLocalOpen}
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
