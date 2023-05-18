import {
  AddOrganizationSubsidiaryMutation,
  useAddOrganizationSubsidiaryMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import {
  GetOrganizationsOptionsDocument,
  GetOrganizationSubsidiariesDocument,
  LinkOrganizationsInput,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';

interface Result {
  onAddOrganizationSubsidiary: (
    input: LinkOrganizationsInput,
  ) => Promise<
    AddOrganizationSubsidiaryMutation['organization_AddSubsidiary'] | null
  >;
}
export const useAddOrganizationSubsidiary = ({
  id,
}: {
  id: string;
}): Result => {
  const [addOrganizationMutation, { loading, error, data }] =
    useAddOrganizationSubsidiaryMutation();
  const handleUpdateCacheAfterAddingOrgSubsidiary = (
    cache: ApolloCache<any>,
    { data: { organization_AddSubsidiary } }: any,
  ) => {
    const data: any | null = client.readQuery({
      query: GetOrganizationSubsidiariesDocument,
      variables: {
        id,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationSubsidiariesDocument,
        variables: {
          id,
        },
        data: {
          organization: {
            subsidiaries: [organization_AddSubsidiary],
          },
        },
      });
      return;
    }

    client.writeQuery({
      query: GetOrganizationSubsidiariesDocument,
      variables: {
        id,
      },
      data: {
        organization: {
          subsidiaries: [organization_AddSubsidiary],
        },
      },
    });
  };

  const handleAddOrganizationSubsidiary: Result['onAddOrganizationSubsidiary'] =
    async (input: LinkOrganizationsInput) => {
      try {
        const response = await addOrganizationMutation({
          variables: { input },
          update: handleUpdateCacheAfterAddingOrgSubsidiary,
        });
        if (response.data?.organization_AddSubsidiary) {
          toast.success('Organization was successfully added!', {
            toastId: `organization-subsidiary-add-success-${response.data.organization_AddSubsidiary.id}`,
          });
        }
        return response.data?.organization_AddSubsidiary ?? null;
      } catch (err) {
        console.error(err);
        toast.error('Something went wrong while adding organization');
        return null;
      }
    };

  return {
    onAddOrganizationSubsidiary: handleAddOrganizationSubsidiary,
  };
};
