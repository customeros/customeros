import { authenticatedClient } from '../../../../../test-resources/nonAdminClients';
import { updateOrganizationInputVariables } from './updateOrganizationInput';
import updateOrganizationGoodExpected from './updateOrganization_good.expected.json';
import {
  CreateOrganizationDocument,
  CreateOrganizationMutation,
  CreateOrganizationMutationVariables,
} from '../../../../organizations/src/graphql/createOrganization.generated';
import {
  UpdateOrganizationDocument,
  UpdateOrganizationMutation,
  UpdateOrganizationMutationVariables,
} from '../../updateOrganization.generated';

describe('graphql suite', () => {
  test('should update organization', async () => {
    let organizationId;
    try {
      const createOrganizationResponse = await authenticatedClient.request<
        CreateOrganizationMutation,
        CreateOrganizationMutationVariables
      >(CreateOrganizationDocument, {
        input: { name: '' },
      });
      organizationId = createOrganizationResponse.organization_Create.id;

      const updatedUpdateOrganizationInput = {
        ...updateOrganizationInputVariables,
      };
      updatedUpdateOrganizationInput.id = organizationId;

      const updateOrganizationResponse = await authenticatedClient.request<
        UpdateOrganizationMutation,
        UpdateOrganizationMutationVariables
      >(UpdateOrganizationDocument, {
        input: updatedUpdateOrganizationInput,
      });

      const updatedExpectedUpdateOrganization = {
        ...updateOrganizationGoodExpected,
      };
      updatedExpectedUpdateOrganization.id = organizationId;

      expect(updateOrganizationResponse.organization_Update).toEqual(
        updatedExpectedUpdateOrganization,
      );
    } catch (e) {
      expect(e).toBeNull();
    }
  });
});
