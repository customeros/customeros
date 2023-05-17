import {
  CreatePhoneCallInteractionEventMutation,
  DataSource,
  GetContactTimelineDocument,
  GetContactTimelineQuery,
  useCreatePhoneCallInteractionEventMutation,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

interface Result {
  onCreatePhoneCallInteractionEvent: (
    input: any,
  ) => Promise<
    CreatePhoneCallInteractionEventMutation['interactionEvent_Create'] | null
  >;
}

const NOW_DATE = new Date().toISOString();

export const useCreatePhoneCallInteractionEvent = ({
  contactId,
}: Props): Result => {
  const [createPhoneCallInteractionEvent, { loading, error, data }] =
    useCreatePhoneCallInteractionEventMutation();

  const handleUpdateCacheAfterAddingPhoneCall = (
    cache: ApolloCache<any>,
    { data: { interactionEvent_Create } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTimelineDocument,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });

    const interactionEvent = {
      ...interactionEvent_Create,
      source: DataSource.Openline,
    };
    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [interactionEvent],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [interactionEvent],
      },
    };

    client.writeQuery({
      query: GetContactTimelineDocument,
      data: newData,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });
  };

  const handleCreatePhoneCallInteractionEvent: Result['onCreatePhoneCallInteractionEvent'] =
    async (input) => {
      try {
        const response = await createPhoneCallInteractionEvent({
          variables: {
            contactId: contactId,
            content: input.content,
            contentType: input.contentType,
            sentBy: input.sentBy,
          },
          update: handleUpdateCacheAfterAddingPhoneCall,
        });
        if (response.data) {
          toast.success('Phone call log added!', {
            toastId: `phone-call-added-${response.data?.interactionEvent_Create.id}`,
          });
        }
        return response.data?.interactionEvent_Create ?? null;
      } catch (err) {
        toast.error('Something went wrong while adding a phone call', {
          toastId: `phone-call-add-error-${contactId}`,
        });
        return null;
      }
    };

  return {
    onCreatePhoneCallInteractionEvent: handleCreatePhoneCallInteractionEvent,
  };
};
