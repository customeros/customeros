import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import { QueryClient } from '@tanstack/react-query';

export const invalidateAccountDetailsQuery = (
  queryClient: QueryClient,
  id: string,
) =>
  queryClient.invalidateQueries(
    useOrganizationAccountDetailsQuery.getKey({ id }),
  );
