import {
  GetOrganizationDetailsQuery,
  UpdateOrganizationIndustryMutation,
  useUpdateOrganizationIndustryMutation,
} from './types';
import {
  GetContactPersonalDetailsWithOrganizationsDocument,
  GetOrganizationDetailsDocument,
  OrganizationUpdateInput,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';

interface Props {
  organizationId: string;
}

interface Result {
  onUpdateOrganizationIndustry: (
    input: Omit<OrganizationUpdateInput, 'id'>,
  ) => Promise<
    UpdateOrganizationIndustryMutation['organization_Update'] | null
  >;
}
export const useUpdateOrganizationIndustry = ({
  organizationId,
}: Props): Result => {
  const [updateOrganizationMutation, { loading, error, data }] =
    useUpdateOrganizationIndustryMutation();

  const handleUpdateCacheAfterUpdatingOrganization = (
    cache: ApolloCache<any>,
    { data: { organization_Update } }: any,
  ) => {
    const data: GetOrganizationDetailsQuery | null = client.readQuery({
      query: GetOrganizationDetailsDocument,
      variables: {
        id: organizationId,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationDetailsDocument,
        data: {
          organization: {
            id: organizationId,
            ...organization_Update,
          },
          variables: { id: organizationId },
        },
      });
    }

    client.writeQuery({
      query: GetOrganizationDetailsDocument,
      data: {
        organization: {
          id: organizationId,
          ...data?.organization,
          industry: organization_Update.industry,
        },
      },
      variables: {
        id: organizationId,
      },
    });
  };

  const handleUpdateOrganizationIndustry: Result['onUpdateOrganizationIndustry'] =
    async (input) => {
      try {
        const response = await updateOrganizationMutation({
          variables: { input: { ...input, id: organizationId } },
          //@ts-expect-error fixme
          update: handleUpdateCacheAfterUpdatingOrganization,
        });
        return response.data?.organization_Update ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateOrganizationIndustry: handleUpdateOrganizationIndustry,
  };
};
