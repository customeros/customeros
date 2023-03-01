import React, { FormEvent } from 'react';
import { useUpdateContactPersonalDetails } from '../../../../hooks/useContact';
import { ContactPersonalDetailsEditForm } from './ContactPersonalDetailsForm';

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
    onUpdateContactPersonalDetails(values).then((value) => {
      if (value) {
        onSetMode('PREVIEW');
      }
    });
  };

  return (
    <ContactPersonalDetailsEditForm
      data={data}
      onSubmit={handleUpdateDetails}
    />
  );
};
