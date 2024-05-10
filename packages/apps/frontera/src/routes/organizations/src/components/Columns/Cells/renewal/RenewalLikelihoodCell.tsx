import { useState } from 'react';

import set from 'lodash/set';
import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { toastError } from '@ui/presentation/Toast';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import { useInfiniteGetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { useBulkUpdateOpportunityRenewalMutation } from '@shared/graphql/bulkUpdateOpportunityRenewal.generated';

import { getLikelihoodColor, getRenewalLikelihoodLabel } from './utils';

interface RenewalLikelihoodCellProps {
  id: string;
  value?: OpportunityRenewalLikelihood | null;
}

export const RenewalLikelihoodCell = ({
  id,
  value,
}: RenewalLikelihoodCellProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [organizationsMeta] = useOrganizationsMeta();
  const colors = value ? getLikelihoodColor(value) : 'text-gray-400';

  const { getOrganization } = organizationsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getOrganization);

  const updateOpportunityRenewal = useBulkUpdateOpportunityRenewalMutation(
    client,
    {
      onMutate: (payload) => {
        queryClient.cancelQueries({ queryKey });

        const { previousEntries } =
          useInfiniteGetOrganizationsQuery.mutateCacheEntry(
            queryClient,
            getOrganization,
          )((old) =>
            produce(old, (draft) => {
              const pageIndex =
                organizationsMeta.getOrganization.pagination.page - 1;

              const content =
                draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content;
              const index = content?.findIndex(
                (item) => item.metadata.id === payload.input.organizationId,
              );

              if (content && index !== undefined && index > -1) {
                set(
                  content[index],
                  'accountDetails.renewalSummary.renewalLikelihood',
                  payload.input.renewalLikelihood,
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
    },
  );

  const handleClick = (value: OpportunityRenewalLikelihood) => {
    updateOpportunityRenewal.mutate({
      input: {
        organizationId: id,
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
        onClick={() => {
          setIsEditing(true);
        }}
        aria-label='edit renewal likelihood'
        icon={<Edit03 className='text-gray-500' />}
        className={cn(
          'rounded-md opacity-0 group-hover/likelihood:opacity-100',
          isEditing && 'opacity-100',
        )}
      />
    </div>
  );
};
