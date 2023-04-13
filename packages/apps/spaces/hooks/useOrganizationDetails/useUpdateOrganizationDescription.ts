import {
  GetOrganizationDetailsQuery,
  UpdateOrganizationDescriptionMutation,
  useUpdateOrganizationDescriptionMutation,
} from './types';
import {
  GetContactPersonalDetailsWithOrganizationsDocument,
  GetOrganizationDetailsDocument,
  OrganizationUpdateInput,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Props {
  organizationId: string;
}

interface Result {
  onUpdateOrganizationDescription: (
    input: Omit<OrganizationUpdateInput, 'id'>,
  ) => Promise<
    UpdateOrganizationDescriptionMutation['organization_Update'] | null
  >;
}
export const useUpdateOrganizationDescription = ({
  organizationId,
}: Props): Result => {
  const [updateOrganizationMutation, { loading, error, data }] =
    useUpdateOrganizationDescriptionMutation();

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
          description: organization_Update.description,
        },
      },
      variables: {
        id: organizationId,
      },
    });
  };

  const handleUpdateOrganizationDescription: Result['onUpdateOrganizationDescription'] =
    async (input) => {
      try {
        const response = await updateOrganizationMutation({
          variables: { input: { ...input, id: organizationId } },
          //@ts-expect-error fixme
          update: handleUpdateCacheAfterUpdatingOrganization,
        });
        return response.data?.organization_Update ?? null;
      } catch (err) {
        toast.error(
          'Something went wrong while updating organization description. Please contact us or try again later',
          {
            toastId: `org-description-${organizationId}-update-error`,
          },
        );
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateOrganizationDescription: handleUpdateOrganizationDescription,
  };
};
