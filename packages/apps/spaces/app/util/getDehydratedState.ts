import { dehydrate } from '@tanstack/react-query';
import getQueryClient from '@shared/util/getQueryClient';
import { getServerGraphQLClient } from './getServerGraphQLClient';

export async function getDehydratedState(hook: any, variables: any) {
  const queryClient = getQueryClient();
  const graphQLClient = getServerGraphQLClient();

  await queryClient.prefetchQuery(
    hook.getKey(variables),
    hook.fetcher(graphQLClient, variables),
  );

  return dehydrate(queryClient);
}
