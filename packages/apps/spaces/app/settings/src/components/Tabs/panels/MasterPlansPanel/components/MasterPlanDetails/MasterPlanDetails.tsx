import { useMemo, useEffect } from 'react';
import { useForm } from 'react-inverted-form';
import { useRouter, useSearchParams } from 'next/navigation';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';
import { useMasterPlansQuery } from '@settings/graphql/masterPlans.generated';
import { useUpdateMasterPlanMutation } from '@settings/graphql/updateMasterPlan.generated';
import { useDuplicateMasterPlanMutation } from '@settings/graphql/duplicateMasterPlan.generated';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { toastError } from '@ui/presentation/Toast';
import { useThrottle } from '@shared/hooks/useThrottle';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { MasterPlanMenu } from './MasterPlanMenu';

interface MasterPlanDetailsProps {
  id: string;
  name: string;
}

type MasterPlanForm = {
  name: string;
};

const formId = 'master-plan-details-form';

export const MasterPlanDetails = ({ id, name }: MasterPlanDetailsProps) => {
  const router = useRouter();
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();
  const defaultValues = useMemo<MasterPlanForm>(() => ({ name }), [name, id]);

  const queryKey = useMasterPlansQuery.getKey();
  const goToPlan = (id: string, options: { retired?: boolean } = {}) => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('planId', id);
    if (options?.retired) {
      params.set('show', 'retired');
    }

    router.push(`/settings?${params.toString()}`);
  };

  const updateMasterPlan = useUpdateMasterPlanMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = useMasterPlansQuery.mutateCacheEntry(
        queryClient,
      )((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const masterPlan = draft?.masterPlans?.find((plan) => plan.id === id);

          if (masterPlan) {
            masterPlan.name = input.name ?? '';

            if (input.retired) {
              masterPlan.retired = input.retired;
            }
          }
        });
      });

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(
        `We couldn't update master plan`,
        'master-plan-details-update',
      );
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });

  const duplicateMasterPlan = useDuplicateMasterPlanMutation(client, {
    onMutate: ({ id }) => {
      queryClient.cancelQueries({ queryKey });

      let tempId = '';

      const { previousEntries } = useMasterPlansQuery.mutateCacheEntry(
        queryClient,
      )((cacheEntry) => {
        return produce(cacheEntry, (draft) => {
          const masterPlan = draft?.masterPlans?.find((plan) => plan.id === id);
          const sameNameCount = draft?.masterPlans?.filter((plan) =>
            plan.name.startsWith(masterPlan?.name ?? ''),
          )?.length;

          tempId = `${masterPlan?.id}-${sameNameCount}`;

          if (masterPlan) {
            draft.masterPlans?.push({
              ...masterPlan,
              id: `${masterPlan.id}-${sameNameCount}`,
              // name: `${masterPlan.name} (copy)`,
            });
          }
        });
      });

      if (tempId) {
        goToPlan(tempId);
      }

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(
        `We couldn't duplicate master plan`,
        'master-plan-duplicate-update',
      );
    },
    onSettled: (data) => {
      queryClient.invalidateQueries({ queryKey });
      if (data) {
        goToPlan(data.masterPlan_Duplicate?.id);
      }
    },
  });

  const handleUpdatePlanName = useThrottle(
    (name: string) => {
      updateMasterPlan.mutate({
        input: {
          id,
          name,
        },
      });
    },
    500,
    [id],
  );

  const { setDefaultValues } = useForm<MasterPlanForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE' && action.payload?.name === 'name') {
        handleUpdatePlanName(action.payload.value);
      }

      return next;
    },
  });

  const handleRetire = () => {
    updateMasterPlan.mutate({
      input: {
        id,
        retired: true,
      },
    });
    goToPlan(id, { retired: true });
  };

  const handleDuplicate = () => {
    duplicateMasterPlan.mutate({ id });
  };

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [name]);

  return (
    <Flex align='center' justify='space-between' mb='2'>
      <FormInput
        name='name'
        formId={formId}
        variant='unstyled'
        borderRadius='unset'
        fontWeight='semibold'
      />
      <MasterPlanMenu onRetire={handleRetire} onDuplicate={handleDuplicate} />
    </Flex>
  );
};
