import { useParams } from 'next/navigation';

import { produce } from 'immer';
import isNil from 'lodash/isNil';
import { useQueryClient } from '@tanstack/react-query';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
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
  const [timelineMeta] = useTimelineMeta();

  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
  const orgQueryKey = useOrganizationQuery.getKey({ id: organizationId });
  const queryKey = useOrganizationOnboardingPlansQuery.getKey({
    organizationId,
  });
  const mutateCacheEntry = useOrganizationOnboardingPlansQuery.mutateCacheEntry(
    queryClient,
    { organizationId },
  );
  const invalidate = () => {
    setTimeout(() => queryClient.invalidateQueries({ queryKey }), 200);
    setTimeout(() => {
      queryClient.invalidateQueries({ queryKey: orgQueryKey });
      queryClient.invalidateQueries({ queryKey: timelineQueryKey });
    }, 2000);
  };

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
            uuid: i?.uuid || '',
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
        `We couldn't update the milestone`,
        'update-org-plan-milestone',
      );
    },
    onSettled: invalidate,
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
      toastError(`We couldn't add the milestone`, 'add-org-plan-milestone');
    },
    onSettled: invalidate,
  });

  return {
    addMilestone,
    updateMilestone,
  };
};
