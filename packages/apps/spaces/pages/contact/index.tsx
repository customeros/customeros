import type { NextPage } from 'next';
import React from 'react';
import Head from 'next/head';
import { ContactList } from '@spaces/contact/contact-list/ContactList';

const ContactsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Contacts</title>
      </Head>
      <ContactList />
    </>
  );
};

export default ContactsPage;
