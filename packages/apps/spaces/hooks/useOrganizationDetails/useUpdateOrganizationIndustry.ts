import {
  GetOrganizationDetailsQuery,
  UpdateOrganizationIndustryMutation,
  useUpdateOrganizationIndustryMutation,
} from './types';
import {
  GetOrganizationDetailsDocument,
  OrganizationUpdateInput,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

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
          update: handleUpdateCacheAfterUpdatingOrganization,
        });
        return response.data?.organization_Update ?? null;
      } catch (err) {
        toast.error(
          'Something went wrong while updating organization industry. Please contact us or try again later',
          {
            toastId: `org-description-${organizationId}-update-error`,
          },
        );
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateOrganizationIndustry: handleUpdateOrganizationIndustry,
  };
};
