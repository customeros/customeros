import React from 'react';
import { ContactDetails } from '@spaces/contact/contact-details/ContactDetails';
import { ContactCommunicationDetails } from '@spaces/contact/contact-communication-details';
import { ContactLocations } from '@spaces/contact/contact-locations';
import { useContact } from '@spaces/hooks/useContact/useGetContact';

export const Details = ({ id }: { id: string }) => {
  const { data, loading, error } = useContact({
    id,
  });

  if (error) {
    return <div>Oops! Something went wrong while loading contact data</div>;
  }

  return (
    <>
      <ContactDetails id={id} data={data} loading={loading}/>
      <ContactCommunicationDetails id={id} data={data} loading={loading} />
      <ContactLocations id={id} data={data} loading={loading} />
    </>
  );
};
