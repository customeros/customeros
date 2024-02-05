import { useParams } from 'next/navigation';

import { produce } from 'immer';
import isNil from 'lodash/isNil';
import { useQueryClient } from '@tanstack/react-query';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  OnboardingPlanMilestoneStatus,
  OnboardingPlanMilestoneItemStatus,
} from '@graphql/types';
import { useOrganizationOnboardingPlansQuery } from '@organization/src/graphql/organizationOnboardingPlans.generated';
import { useAddOnboardingPlanMilestoneMutation } from '@organization/src/graphql/addOnboardingPlanMilestone.generated';
import { useUpdateOnboardingPlanMilestoneMutation } from '@organization/src/graphql/updateOnboardingPlanMilestone.generated';

import { PlanDatum } from '../OnboardingPlans/types';

interface UseMilestoneMutationsOptions {
  plan?: PlanDatum;
}

export const useMilestoneMutations = (
  options: UseMilestoneMutationsOptions = {},
) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const organizationId = (useParams()?.id ?? '') as string;

  const queryKey = useOrganizationOnboardingPlansQuery.getKey({
    organizationId,
  });
  const mutateCacheEntry = useOrganizationOnboardingPlansQuery.mutateCacheEntry(
    queryClient,
    { organizationId },
  );
  const invalidateQuery = () => queryClient.invalidateQueries({ queryKey });

  const updateMilestone = useUpdateOnboardingPlanMilestoneMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateCacheEntry((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const task = draft?.organizationPlansForOrganization?.find(
            (plan) => plan.id === input.organizationPlanId,
          );
          if (!task) return;

          const milestone = task.milestones.find(
            (milestone) => milestone.id === input.id,
          );
          if (!milestone) return;

          if (input.statusDetails) {
            milestone.statusDetails = input.statusDetails;
          }
          if (!isNil(input.retired)) {
            milestone.retired = input.retired;
          }

          milestone.items = (input?.items ?? []).map((i) => ({
            status: i?.status || OnboardingPlanMilestoneItemStatus.NotDone,
            text: i?.text || '',
            updatedAt: i?.updatedAt ?? '',
          }));
        });
      });

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        mutateCacheEntry(() => context.previousEntries);
      }
      toastError(
        `We could'nt update the milestone`,
        'update-org-plan-milestone',
      );
    },
    onSettled: invalidateQuery,
  });

  const addMilestone = useAddOnboardingPlanMilestoneMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateCacheEntry((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const plan = draft?.organizationPlansForOrganization?.find(
            (plan) => plan.id === input.organizationPlanId,
          );
          if (!plan) return;

          plan.milestones.push({
            ...input,
            id: 'temp',
            name: input.name ?? '',
            items: [],
            retired: false,
            statusDetails: {
              status: OnboardingPlanMilestoneStatus.NotStarted,
              text: '',
              updatedAt: new Date().toISOString(),
            },
          });
        });
      });

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        mutateCacheEntry(() => context.previousEntries);
      }
      toastError(`We could'nt add the milestone`, 'add-org-plan-milestone');
    },
    onSettled: invalidateQuery,
  });

  return {
    addMilestone,
    updateMilestone,
  };
};
