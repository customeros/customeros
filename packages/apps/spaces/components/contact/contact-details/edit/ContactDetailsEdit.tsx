import React, { FormEvent } from 'react';
import { useUpdateContactPersonalDetails } from '../../../../hooks/useContact';
import { ContactPersonalDetailsEditForm } from './ContactDetailsForm';

export const ContactPersonalDetailsEdit = ({
  data,
  onSetMode,
}: {
  data: any;
  onSetMode: any;
}) => {
  const { onUpdateContactPersonalDetails } = useUpdateContactPersonalDetails({
    contactId: data.id,
  });

  const handleUpdateDetails = (values: any) => {
    onUpdateContactPersonalDetails(values).then(() => {
      onSetMode('PREVIEW');
    });
  };

  return (
    <ContactPersonalDetailsEditForm
      data={data}
      onSubmit={handleUpdateDetails}
    />
  );
};
