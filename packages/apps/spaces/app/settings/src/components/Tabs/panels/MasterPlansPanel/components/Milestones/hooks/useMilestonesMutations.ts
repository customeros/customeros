import { useSearchParams } from 'next/navigation';

import set from 'lodash/set';
import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useMasterPlansQuery } from '@settings/graphql/masterPlans.generated';
import { useCreateMilestoneMutation } from '@settings/graphql/createMilestone.generated';
import { useUpdateMilestoneMutation } from '@settings/graphql/updateMilestone.generated';
import { useUpdateMilestonesMutation } from '@settings/graphql/updateMilestones.generated';
import { useDuplicateMilestoneMutation } from '@settings/graphql/duplicateMilestone.generated';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const useMilestonesMutations = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const planId = useSearchParams()?.get('planId') ?? '';

  const queryKey = useMasterPlansQuery.getKey();
  const mutateMasterPlansCache =
    useMasterPlansQuery.mutateCacheEntry(queryClient);
  const invalidateMasterPlans = useDebounce(
    () => queryClient.invalidateQueries({ queryKey }),
    500,
  );

  const createMilestone = useCreateMilestoneMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateMasterPlansCache((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const idx = draft?.masterPlans?.findIndex(
            (mp) => mp.id === input?.masterPlanId,
          );

          if (idx === -1) return;
          const allCachedMilestones = draft?.masterPlans?.[idx]?.milestones;
          const lastMilestone =
            allCachedMilestones?.[allCachedMilestones?.length - 1];
          const currentMilestonesCount = allCachedMilestones?.length ?? 0;

          draft?.masterPlans?.[idx]?.milestones.push({
            ...input,
            retired: false,
            name: input?.name ?? '',
            order: lastMilestone?.order ? lastMilestone?.order + 1 : 0,
            id: `${currentMilestonesCount + 1}`,
          });
        });
      });

      return { previousEntries };
    },
    onError: (e, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(`We couldn't create the milestone`, 'create-milestone-error');
    },
    onSettled: () => {
      invalidateMasterPlans();
    },
  });

  const updateMilestone = useUpdateMilestoneMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateMasterPlansCache((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const masterPlanIdx = draft?.masterPlans?.findIndex(
            (mp) => mp.id === input?.masterPlanId,
          );

          if (masterPlanIdx === -1) return;
          const milestoneIdx = draft?.masterPlans?.[
            masterPlanIdx
          ]?.milestones?.findIndex((m) => m.id === input?.id);

          if (milestoneIdx === -1) return;

          if (input.retired) {
            draft?.masterPlans?.[masterPlanIdx]?.milestones?.splice(
              milestoneIdx,
              1,
            );

            return;
          }

          set(
            draft,
            `masterPlans.${masterPlanIdx}.milestones.${milestoneIdx}`,
            { ...input },
          );
        });
      });

      return { previousEntries };
    },
    onError: (e, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(`We couldn't update the milestone`, 'update-milestone-error');
    },
    onSettled: () => {
      invalidateMasterPlans();
    },
  });

  const updateMilestones = useUpdateMilestonesMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateMasterPlansCache((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const masterPlanIdx = draft?.masterPlans?.findIndex(
            (mp) => mp.id === planId,
          );

          if (masterPlanIdx === -1) return;

          set(draft, `masterPlans.${masterPlanIdx}.milestones`, input);
        });
      });

      return { previousEntries };
    },
    onError: (e, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(
        `We couldn't update the milestones`,
        'update-bulk-milestones-error',
      );
    },
    onSettled: () => {
      invalidateMasterPlans();
    },
  });

  const duplicateMilestone = useDuplicateMilestoneMutation(client, {
    onMutate: ({ id, masterPlanId }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateMasterPlansCache((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const masterPlan = draft?.masterPlans?.find(
            (mp) => mp.id === masterPlanId,
          );
          if (!masterPlan) return;

          const milestone = masterPlan?.milestones?.find((m) => m.id === id);
          if (!milestone) return;

          const currentMilestonesCount = masterPlan?.milestones?.length ?? 0;
          masterPlan.milestones.push({
            ...milestone,
            id: `${currentMilestonesCount + 1}`,
            order: currentMilestonesCount + 1,
          });
        });
      });

      return { previousEntries };
    },
    onError: (e, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(
        `We couldn't duplicate the milestone`,
        'duplicate-milestones-error',
      );
    },
    onSettled: () => {
      invalidateMasterPlans();
    },
  });

  return {
    createMilestone,
    updateMilestone,
    updateMilestones,
    duplicateMilestone,
  };
};
