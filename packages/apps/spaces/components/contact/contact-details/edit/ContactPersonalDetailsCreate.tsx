import React from 'react';
import { useCreateContact, ContactInput } from '../../../../hooks/useContact';
import { ContactPersonalDetailsEditForm } from './ContactPersonalDetailsForm';
import { useRouter } from 'next/router';

export const ContactPersonalDetailsCreate = () => {
  const router = useRouter();
  const { onCreateContact } = useCreateContact();

  const handleCreateDetails = ({
    ownerFullName,
    ...values
  }: // fixme ownerFullName is connected to search owner component
  ContactInput & { ownerFullName: string }) => {
    onCreateContact(values).then((value) => {
      if (value?.id) {
        router.push(`/contact/${value.id}`);
      }
    });
  };

  return (
    <ContactPersonalDetailsEditForm
      data={{}}
      onSubmit={handleCreateDetails}
      mode={'CREATE'}
    />
  );
};
