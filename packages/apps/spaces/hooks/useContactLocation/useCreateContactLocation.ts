import {
  AddLocationToContactMutation,
  GetContactLocationsQuery,
  useAddLocationToContactMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import { GetContactLocationsDocument } from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';

interface Props {
  contactId: string;
}

interface Result {
  saving: boolean;
  onCreateContactLocation: () => Promise<
    AddLocationToContactMutation['contact_AddNewLocation'] | null
  >;
}

export const useCreateContactLocation = ({ contactId }: Props): Result => {
  const [createContactLocationMutation, { loading }] =
    useAddLocationToContactMutation();

  const handleUpdateCacheAfterAddingLocation = (
    cache: ApolloCache<any>,
    { data: { contact_AddNewLocation } }: any,
  ) => {
    const data: GetContactLocationsQuery | null = client.readQuery({
      query: GetContactLocationsDocument,
      variables: {
        id: contactId,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetContactLocationsDocument,
        variables: {
          id: contactId,
        },
        data: {
          contact: {
            id: contactId,
            locations: [contact_AddNewLocation],
          },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        locations: [
          ...(data.contact?.locations || []),
          { ...contact_AddNewLocation },
        ],
      },
    };
    client.writeQuery({
      query: GetContactLocationsDocument,
      data: newData,
      variables: {
        id: contactId,
      },
    });
  };

  const handleCreateContactLocation: Result['onCreateContactLocation'] =
    async () => {
      try {
        const response = await createContactLocationMutation({
          variables: { contactId },

          update: handleUpdateCacheAfterAddingLocation,
        });
        if (response.data) {
          toast.success('Location added!', {
            toastId: `location-added-${response.data?.contact_AddNewLocation.id}`,
          });
        }
        return response.data?.contact_AddNewLocation ?? null;
      } catch (err) {
        toast.error('Something went wrong while adding location', {
          toastId: `Location-add-error-${contactId}`,
        });
        return null;
      }
    };

  return {
    saving: loading,
    onCreateContactLocation: handleCreateContactLocation,
  };
};
