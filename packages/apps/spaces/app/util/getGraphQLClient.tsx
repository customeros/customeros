import { GraphQLClient } from 'graphql-request';

export const getGraphQLClient = () => {
  return new GraphQLClient(`${process.env.CUSTOMER_OS_API_PATH}/query`, {
    headers: {
      // 'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY as string,
      // 'X-Openline-USERNAME': 'development@openline.ai',

      'X-Openline-API-KEY': 'dd9d2474-b4a9-4799-b96f-73cd0a2917e4',
      'X-Openline-USERNAME': 'acalinica@openline.ai',
    },
  });
};
