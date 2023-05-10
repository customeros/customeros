import React from 'react';
import { useContactCommunicationChannelsDetails } from '@spaces/hooks/useContact/useContactCommunicationChannelsDetails';
import {
  useAddEmailToContactEmail,
  useRemoveEmailFromContactEmail,
  useUpdateContactEmail,
} from '@spaces/hooks/useContactEmail';
import {
  useCreateContactPhoneNumber,
  useRemovePhoneNumberFromContact,
  useUpdateContactPhoneNumber,
} from '@spaces/hooks/useContactPhoneNumber';

import { useRecoilValue } from 'recoil';
import { contactDetailsEdit } from '../../../state';
import { CommunicationDetails } from '@spaces/molecules/communication-details';

export const ContactCommunicationDetails = ({ id }: { id: string }) => {
  const { isEditMode } = useRecoilValue(contactDetailsEdit);

  const { data, loading, error } = useContactCommunicationChannelsDetails({
    id,
  });

  const { onAddEmailToContact } = useAddEmailToContactEmail({ contactId: id });

  const { onRemoveEmailFromContact } = useRemoveEmailFromContactEmail({
    contactId: id,
  });
  const { onUpdateContactEmail } = useUpdateContactEmail({
    contactId: id,
  });

  const { onCreateContactPhoneNumber } = useCreateContactPhoneNumber({
    contactId: id,
  });
  const { onUpdateContactPhoneNumber } = useUpdateContactPhoneNumber({
    contactId: id,
  });
  const { onRemovePhoneNumberFromContact } = useRemovePhoneNumberFromContact({
    contactId: id,
  });

  return (
    <CommunicationDetails
      id={id}
      onAddEmail={(input: any) => onAddEmailToContact(input)}
      onAddPhoneNumber={(input: any) => onCreateContactPhoneNumber(input)}
      onRemoveEmail={(input: any) => onRemoveEmailFromContact(input)}
      onRemovePhoneNumber={(input: any) =>
        onRemovePhoneNumberFromContact(input)
      }
      onUpdateEmail={(input: any) => onUpdateContactEmail(input)}
      onUpdatePhoneNumber={(input: any) => onUpdateContactPhoneNumber(input)}
      // @ts-expect-error fixme
      data={data}
      loading={loading}
      isEditMode={isEditMode}
    />
  );
};
