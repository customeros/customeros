import { dehydrate } from '@tanstack/react-query';
import getQueryClient from '@shared/util/getQueryClient';
import { getServerGraphQLClient } from './getServerGraphQLClient';

// eslint-disable-next-line
export async function getDehydratedState(hook: any, variables?: any) {
  const queryClient = getQueryClient();
  const graphQLClient = getServerGraphQLClient();

  try {
    await queryClient.prefetchQuery(
      hook.getKey(variables),
      hook.fetcher(graphQLClient, variables),
    );
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('getDehydratedState: ', error);
  }

  return dehydrate(queryClient);
}
