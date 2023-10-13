import { TenantNameDocument } from '../../tenantName.generated';
import expectedGoodTenantName from './tenantName_good.expected.json';
import expectedBadTenantName from './tenantName_bad.expected.json';
import expectedUnauthenticatedUser from './unauthenticatedUser.expected.json';
import {
  authenticatedClient,
  unauthenticatedClient,
} from '../../../../../test-resources/nonAdminClients';

describe('graphql suite', () => {
  test('should match tenant name', async () => {
    try {
      const res = await authenticatedClient.request(TenantNameDocument);
      expect(res).toEqual(expectedGoodTenantName);
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});

describe('graphql suite', () => {
  test('should not match tenant name', async () => {
    try {
      const res = await authenticatedClient.request(TenantNameDocument);
      expect(res).not.toEqual(expectedBadTenantName);
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});

describe('graphql suite', () => {
  test('should return error for unauthenticated user', async () => {
    try {
      await unauthenticatedClient.request(TenantNameDocument);
    } catch (e: any) {
      const filteredReceived = Object.fromEntries(
        Object.entries(e.response).filter(([key]) => key !== 'headers'),
      );
      expect(filteredReceived).toEqual(expectedUnauthenticatedUser);

    }
  });
});
