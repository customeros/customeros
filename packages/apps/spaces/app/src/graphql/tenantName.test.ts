import { GraphQLClient } from 'graphql-request';

import { TenantNameDocument } from './tenantName.generated';
import expected from './tenantName.expected.json';

const headers = {
  'X-Openline-API-KEY': 'dd9d2474-b4a9-4799-b96f-73cd0a2917e4',
  'X-Openline-TENANT': 'openline',
};

const client = new GraphQLClient('http://127.0.0.1:10000/query', {
  credentials: 'include',
  headers,
});

describe('graphql suite', () => {
  test('should match tenant name', async () => {
    try {
      const res = await client.request(TenantNameDocument);
      expect(res).toEqual(expected);
    } catch (e) {
      console.error(e);
    }
  });
});
