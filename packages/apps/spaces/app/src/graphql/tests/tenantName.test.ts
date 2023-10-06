import { TenantNameDocument } from '../tenantName.generated';
import expected from './tenantName.expected.json';
import { client } from '../../../../test.resources';

describe('graphql suite', () => {
  test('should match tenant name', async () => {
    try {
      const res = await client.request(TenantNameDocument);
      expect(res).toEqual(expected);
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});
