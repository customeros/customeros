import { GraphQLClient } from 'graphql-request';

export const getGraphQLClient = () => {
  return new GraphQLClient(`${process.env.CUSTOMER_OS_API_PATH}/query`, {
    headers: {
      'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY as string,
      'X-Openline-USERNAME': 'development@openline.ai',
    },
  });
};
