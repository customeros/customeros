import { authenticatedClient } from '../../../../../../test-resources/nonAdminClients';
import {
  CreateOrganizationDocument,
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from '../../createOrganization.generated';

describe('graphql suite', () => {
  test('should create organization', async () => {
    try {
      const createOrganizationResponse = await authenticatedClient.request<
        CreateOrganizationMutation,
        CreateOrganizationMutationVariables
      >(CreateOrganizationDocument, {
        input: { name: '' },
      });
      expect(createOrganizationResponse.organization_Create.id.length).toEqual(36);
      expect(createOrganizationResponse.organization_Create.name).toEqual('');
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});
