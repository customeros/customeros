import { authenticatedClient } from '../../../../../../test-resources/nonAdminClients';
import { createOrganizationInputVariables } from './createOrganizationInput';
import {
  CreateOrganizationDocument,
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from '../../createOrganization.generated';
import gql from 'graphql-tag';
//

describe('graphql suite', () => {
  test('should create organization', async () => {
    try {
      // const res = await authenticatedClient.request(mutation, {
      //   input: createOrganizationInputVariables,
      // });
      const res = await authenticatedClient.request<
        CreateOrganizationMutation,
        CreateOrganizationMutationVariables
      >(CreateOrganizationDocument, {
        input: { name: '' },
      });
      expect(res.organization_Create.id.length).toEqual(36);
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});
