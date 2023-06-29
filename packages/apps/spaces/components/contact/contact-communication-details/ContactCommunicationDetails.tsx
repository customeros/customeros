import React from 'react';
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

export const ContactCommunicationDetails = ({
  id,
  data,
  loading,
}: {
  id: string;
  data: any;
  loading: any;
}) => {
  const { isEditMode } = useRecoilValue(contactDetailsEdit);

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
      onAddEmail={(input: any) => onAddEmailToContact(input)}
      onAddPhoneNumber={(input: any) => onCreateContactPhoneNumber(input)}
      onRemoveEmail={(input: any) => onRemoveEmailFromContact(input)}
      onRemovePhoneNumber={(input: any) =>
        onRemovePhoneNumberFromContact(input)
      }
      onUpdateEmail={(input: any) => onUpdateContactEmail(input)}
      onUpdatePhoneNumber={(input: any) => onUpdateContactPhoneNumber(input)}
      data={data}
      loading={loading}
      isEditMode={isEditMode}
    />
  );
};
