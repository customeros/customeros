import {
  GetOrganizationDetailsQuery,
  UpdateOrganizationNameMutation,
  useUpdateOrganizationNameMutation,
} from './types';
import {
  GetContactPersonalDetailsWithOrganizationsDocument,
  GetOrganizationDetailsDocument,
  OrganizationUpdateInput,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';
import { toast } from 'react-toastify';

interface Props {
  organizationId: string;
}

interface Result {
  onUpdateOrganizationName: (
    input: Omit<OrganizationUpdateInput, 'id'>,
  ) => Promise<UpdateOrganizationNameMutation['organization_Update'] | null>;
}
export const useUpdateOrganizationName = ({
  organizationId,
}: Props): Result => {
  const [updateOrganizationMutation, { loading, error, data }] =
    useUpdateOrganizationNameMutation();

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
          name: organization_Update.name,
        },
      },
      variables: {
        id: organizationId,
      },
    });
  };

  const handleUpdateOrganizationName: Result['onUpdateOrganizationName'] =
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
          'Something went wrong while updating organization name. Please contact us or try again later',
          {
            toastId: `org-name-${organizationId}-update-error`,
          },
        );
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateOrganizationName: handleUpdateOrganizationName,
  };
};
