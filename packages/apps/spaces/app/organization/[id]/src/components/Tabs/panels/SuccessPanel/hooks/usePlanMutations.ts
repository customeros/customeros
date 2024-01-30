import { produce } from 'immer';
import isNil from 'lodash/isNil';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCreateOnboardingPlanMutation } from '@organization/src/graphql/createOnboardingPlan.generated';
import { useUpdateOnboardingPlanMutation } from '@organization/src/graphql/updateOnboardingPlan.generated';
import { useOrganizationOnboardingPlansQuery } from '@organization/src/graphql/organizationOnboardingPlans.generated';

interface UsePlanMutationsOptions {
  organizationId: string;
}

export const usePlanMutations = ({
  organizationId,
}: UsePlanMutationsOptions) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const queryKey = useOrganizationOnboardingPlansQuery.getKey({
    organizationId,
  });
  const mutateCacheEntry = useOrganizationOnboardingPlansQuery.mutateCacheEntry(
    queryClient,
    { organizationId },
  );
  const invalidateQuery = useDebounce(
    () => queryClient.invalidateQueries({ queryKey }),
    500,
  );

  const createOnboardingPlan = useCreateOnboardingPlanMutation(client, {
    onMutate: ({ input: { name } }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateCacheEntry((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          draft?.organizationPlansForOrganization?.unshift({
            masterPlanId: 'temp',
            id: 'temp',
            name: name ?? 'Unnamed plan',
            milestones: [],
            retired: false,
            statusDetails: {
              status: 'NOT_STARTED',
              text: 'Not started',
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
      toastError(`We could'nt create the plan`, 'create-org-onboarding-plan');
    },
    onSettled: invalidateQuery,
  });

  const updateOnboardingPlan = useUpdateOnboardingPlanMutation(client, {
    onMutate: ({ input: { id, retired } }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = mutateCacheEntry((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const plan = draft?.organizationPlansForOrganization?.find(
            (plan) => plan.id === id,
          );
          if (plan) {
            if (isNil(retired)) return;
            plan.retired = retired;
          }
        });
      });

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        mutateCacheEntry(() => context.previousEntries);
      }
      toastError(`We could'nt update the plan`, 'update-org-onboarding-plan');
    },
    onSettled: invalidateQuery,
  });

  return {
    createOnboardingPlan,
    updateOnboardingPlan,
  };
};
