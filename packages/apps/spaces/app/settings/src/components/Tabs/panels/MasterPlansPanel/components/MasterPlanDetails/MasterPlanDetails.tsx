import { useMemo, useEffect } from 'react';
import { useForm } from 'react-inverted-form';
import { useRouter, useSearchParams } from 'next/navigation';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';
import { useMasterPlansQuery } from '@settings/graphql/masterPlans.generated';
import { useUpdateMasterPlanMutation } from '@settings/graphql/updateMasterPlan.generated';
import { useDuplicateMasterPlanMutation } from '@settings/graphql/duplicateMasterPlan.generated';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { FormInput } from '@ui/form/Input';
import { toastError } from '@ui/presentation/Toast';
import { useThrottle } from '@shared/hooks/useThrottle';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { MasterPlanMenu } from './MasterPlanMenu';

interface MasterPlanDetailsProps {
  id: string;
  name: string;
  isRetired?: boolean;
}

type MasterPlanForm = {
  name: string;
};

const formId = 'master-plan-details-form';

export const MasterPlanDetails = ({
  id,
  name,
  isRetired,
}: MasterPlanDetailsProps) => {
  const router = useRouter();
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();
  const defaultValues = useMemo<MasterPlanForm>(() => ({ name }), [name, id]);

  const queryKey = useMasterPlansQuery.getKey();
  const goToPlan = (id: string, options: { retired?: boolean } = {}) => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('planId', id);
    params.set('show', options?.retired ? 'retired' : 'active');

    router.push(`/settings?${params.toString()}`);
  };

  const updateMasterPlan = useUpdateMasterPlanMutation(client, {
    onError: (_, __, context) => {
      toastError(
        `We couldn't update master plan`,
        'master-plan-details-update',
      );
    },
    onSettled: (_, __, { input }) => {
      queryClient.invalidateQueries({ queryKey });
      goToPlan(input.id, { retired: input.retired ?? false });
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

  const isLoading = updateMasterPlan.isPending || duplicateMasterPlan.isPending;

  const handleRetire = () => {
    updateMasterPlan.mutate({
      input: {
        id,
        retired: true,
      },
    });
  };

  const handleDuplicate = () => {
    duplicateMasterPlan.mutate({ id });
  };

  const handleReactivate = () => {
    updateMasterPlan.mutate({
      input: {
        id,
        retired: false,
      },
    });
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
        defaultValue={name}
        borderRadius='unset'
        fontWeight='semibold'
      />
      {isRetired ? (
        <Button
          size='xs'
          variant='outline'
          isLoading={isLoading}
          onClick={handleReactivate}
        >
          Reactivate
        </Button>
      ) : (
        <MasterPlanMenu
          isLoading={isLoading}
          onRetire={handleRetire}
          onDuplicate={handleDuplicate}
        />
      )}
    </Flex>
  );
};
