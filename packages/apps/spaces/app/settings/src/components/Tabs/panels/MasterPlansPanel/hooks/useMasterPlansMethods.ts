import { useRouter, useSearchParams } from 'next/navigation';

import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useMasterPlansQuery } from '@settings/graphql/masterPlans.generated';
import { useCreateMasterPlanMutation } from '@settings/graphql/createMasterPlan.generated';
import { useCreateDefaultMasterPlanMutation } from '@settings/graphql/createDefaultMasterPlan.generated';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const useMasterPlansMethods = () => {
  const router = useRouter();
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();

  const queryKey = useMasterPlansQuery.getKey();
  const invalidateQuery = useDebounce(
    () => queryClient.invalidateQueries({ queryKey }),
    500,
  );

  const goToNewPlan = (id = '') => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');

    params.set('planId', id);
    router.push(`/settings?${params.toString()}`);
  };

  const createMasterPlan = useCreateMasterPlanMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.getQueryData(queryKey);

      let tempId = '';

      const { previousEntries } = useMasterPlansQuery.mutateCacheEntry(
        queryClient,
      )((cacheEntry) => {
        tempId = `${cacheEntry?.masterPlans?.length + 1}`;

        return produce(cacheEntry, (draft) => {
          draft?.masterPlans?.push({
            id: tempId,
            name: input.name ?? 'New Master Plan',
            retired: false,
            milestones: [],
            retiredMilestones: [],
          });
        });
      });

      goToNewPlan(tempId);

      return { previousEntries, queryKey };
    },
    onError: (_, __, context) => {
      if (context?.queryKey && context?.previousEntries) {
        queryClient.setQueryData(context?.queryKey, context?.previousEntries);
      }
      toastError(`We couldn't create the master plan`, 'create-master-plan');
    },
    onSettled: (data) => {
      if (data) {
        const newPlanId = data.masterPlan_Create?.id;
        goToNewPlan(newPlanId);
      }
      invalidateQuery();
    },
  });

  const createDefaultMasterPlan = useCreateDefaultMasterPlanMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.getQueryData(queryKey);

      let tempId = '';

      const { previousEntries } = useMasterPlansQuery.mutateCacheEntry(
        queryClient,
      )((cacheEntry) => {
        tempId = `${cacheEntry?.masterPlans?.length + 1}`;

        return produce(cacheEntry, (draft) => {
          draft?.masterPlans?.push({
            id: tempId,
            name: input.name ?? 'New Master Plan',
            retired: false,
            milestones: [],
            retiredMilestones: [],
          });
        });
      });

      goToNewPlan(tempId);

      return { previousEntries, queryKey };
    },
    onError: (_, __, context) => {
      if (context?.queryKey && context?.previousEntries) {
        queryClient.setQueryData(context?.queryKey, context?.previousEntries);
      }
      toastError(
        `We couldn't create the master plan`,
        'create-master-default-plan',
      );
    },
    onSettled: (data) => {
      if (data) {
        const newPlanId = data.masterPlan_CreateDefault?.id;
        goToNewPlan(newPlanId);
      }
      invalidateQuery();
    },
  });

  const handleCreateFromScratch = () => {
    createMasterPlan.mutate({
      input: {
        name: 'New Master Plan',
      },
    });
  };
  const handleCreateDefault = () => {
    createDefaultMasterPlan.mutate({
      input: {
        name: 'Default Master Plan',
      },
    });
  };

  return {
    handleCreateFromScratch,
    handleCreateDefault,
    isPending: createMasterPlan.isPending || createDefaultMasterPlan.isPending,
  };
};
