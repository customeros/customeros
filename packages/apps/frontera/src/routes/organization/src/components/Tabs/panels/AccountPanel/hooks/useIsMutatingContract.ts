import { useIsMutating } from '@tanstack/react-query';

import { useCreateServiceMutation } from '@organization/graphql/createService.generated';
import { useUpdateServiceMutation } from '@organization/graphql/updateService.generated';
import { useUpdateOpportunityRenewalMutation } from '@organization/graphql/updateOpportunityRenewal.generated';

export const useIsMutatingContract = () => {
  const isMutatingLikelihood = useIsMutating({
    mutationKey: useUpdateOpportunityRenewalMutation.getKey(),
  });
  const isMutatingOpportunityRenewal = useIsMutating({
    mutationKey: useUpdateOpportunityRenewalMutation.getKey(),
  });
  const isMutatingCreateService = useIsMutating({
    mutationKey: useCreateServiceMutation.getKey(),
  });
  const isMutatingUpdateService = useIsMutating({
    mutationKey: useUpdateServiceMutation.getKey(),
  });

  return (
    isMutatingLikelihood +
    isMutatingOpportunityRenewal +
    isMutatingCreateService +
    isMutatingUpdateService
  );
};
