import { useState } from 'react';

import set from 'lodash/set';
import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';
import { useInfiniteGetRenewalsQuery } from '@renewals/graphql/getRenewals.generated';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { toastError } from '@ui/presentation/Toast';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useRenewalsMeta } from '@shared/state/RenewalsMeta.atom';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import { useInfiniteGetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { useUpdateOpportunityRenewalMutation } from '@organization/graphql/updateOpportunityRenewal.generated';

import { getLikelihoodColor, getRenewalLikelihoodLabel } from './utils';

interface RenewalLikelihoodCellProps {
  id: string;
  opportunityId: string;
  value?: OpportunityRenewalLikelihood | null;
}

export const RenewalLikelihoodCell = ({
  id,
  value,
  opportunityId,
}: RenewalLikelihoodCellProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [renewalsMeta] = useRenewalsMeta();

  const colors = value ? getLikelihoodColor(value) : 'text-gray-400';

  const { getRenewals } = renewalsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getRenewals);

  const updateOpportunityRenewal = useUpdateOpportunityRenewalMutation(client, {
    onMutate: (payload) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = useInfiniteGetRenewalsQuery.mutateCacheEntry(
        queryClient,
        getRenewals,
      )((old) =>
        produce(old, (draft) => {
          const pageIndex = getRenewals.pagination.page - 1;

          const content =
            draft?.pages?.[pageIndex]?.dashboardView_Renewals?.content;
          const index = content?.findIndex(
            (item) => item.organization.metadata.id === id,
          );

          if (content && index !== undefined && index > -1) {
            set(
              content[index],
              'opportunity.renewalLikelihood',
              payload.input.renewalLikelihood,
            );
            set(
              content[index],
              'opportunity.renewalAdjustedRate',
              payload.input.renewalAdjustedRate,
            );
          }
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      toastError(
        `We couldn't update the renewal likelihood`,
        'renewal-likelihood-update-error',
      );

      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
    },
    onSettled: () => {
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
      }, 500);
    },
  });

  const handleClick = (value: OpportunityRenewalLikelihood) => {
    updateOpportunityRenewal.mutate({
      input: {
        opportunityId,
        renewalLikelihood: value,
        renewalAdjustedRate: (() => {
          switch (value) {
            case OpportunityRenewalLikelihood.HighRenewal:
              return 100;
            case OpportunityRenewalLikelihood.MediumRenewal:
              return 50;
            case OpportunityRenewalLikelihood.LowRenewal:
              return 25;
            default:
              return 50;
          }
        })(),
      },
    });
  };

  return (
    <div className='flex gap-1 items-center group/likelihood'>
      <Menu open={isEditing} onOpenChange={setIsEditing}>
        <MenuButton asChild disabled>
          <span
            className={cn('cursor-default', colors)}
            onDoubleClick={() => setIsEditing(true)}
          >
            {value ? getRenewalLikelihoodLabel(value) : 'Unknown'}
          </span>
        </MenuButton>
        <MenuList align='center'>
          <MenuItem
            onClick={() =>
              handleClick(OpportunityRenewalLikelihood.HighRenewal)
            }
          >
            High
          </MenuItem>
          <MenuItem
            onClick={() =>
              handleClick(OpportunityRenewalLikelihood.MediumRenewal)
            }
          >
            Medium
          </MenuItem>
          <MenuItem
            onClick={() => handleClick(OpportunityRenewalLikelihood.LowRenewal)}
          >
            Low
          </MenuItem>
        </MenuList>
      </Menu>

      <IconButton
        size='xxs'
        variant='ghost'
        aria-label='edit renewal likelihood'
        icon={<Edit03 className='text-gray-500' />}
        onClick={() => {
          setIsEditing(true);
        }}
        className={cn(
          'rounded-md opacity-0 group-hover/likelihood:opacity-100',
          isEditing && 'opacity-100',
        )}
      />
    </div>
  );
};
