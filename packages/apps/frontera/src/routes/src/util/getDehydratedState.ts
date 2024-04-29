import { dehydrate } from '@tanstack/react-query';

import getQueryClient from '@shared/util/getQueryClient';

import { getServerGraphQLClient } from './getServerGraphQLClient';

export async function getDehydratedState(
  // eslint-disable-next-line
  hook: any,
  // eslint-disable-next-line
  options?: { variables?: any; fetcher?: any },
) {
  const queryClient = getQueryClient();
  const graphQLClient = getServerGraphQLClient();

  try {
    await queryClient.prefetchQuery({
      queryKey: hook.getKey(options?.variables),
      queryFn: (options?.fetcher ? options.fetcher : hook.fetcher)(
        graphQLClient,
        options?.variables,
      ),
    });
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('getDehydratedState: ', error);
  }

  return dehydrate(queryClient);
}
