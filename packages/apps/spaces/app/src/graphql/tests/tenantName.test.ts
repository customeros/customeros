import { TenantNameDocument } from '../tenantName.generated';
import expectedGoodTenantName from './tenantName_good.expected.json';
import expectedBadTenantName from './tenantName_bad.expected.json';
import expectedUnauthenticatedUser from './unauthenticatedUser.expected.json';
import {
  authenticated_client,
  unauthenticated_client,
} from '../../../../test.resources';

describe('graphql suite', () => {
  test('should match tenant name', async () => {
    try {
      const res = await authenticated_client.request(TenantNameDocument);
      expect(res).toEqual(expectedGoodTenantName);
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});

describe('graphql suite', () => {
  test('should not match tenant name', async () => {
    try {
      const res = await authenticated_client.request(TenantNameDocument);
      expect(res).not.toEqual(expectedBadTenantName);
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});

describe('graphql suite', () => {
  test('should return error for unauthenticated user', async () => {
    try {
      await unauthenticated_client.request(TenantNameDocument);
      // expect(res).not.toEqual(unauthenticatedUser);
    } catch (e: any) {
      const filteredReceived = Object.fromEntries(
        Object.entries(e.response).filter(([key]) => key !== 'headers'),
      );
      expect(filteredReceived).toEqual(expectedUnauthenticatedUser);

    }
  });
});
