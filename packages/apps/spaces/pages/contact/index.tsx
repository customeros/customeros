import type { NextPage } from 'next';
import React from 'react';
import Head from 'next/head';
import { ContactList } from '@spaces/contact/contact-list/ContactList';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';

const ContactsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Contacts</title>
      </Head>
      <PageContentLayout>
        <ContactList />
      </PageContentLayout>
    </>
  );
};

export default ContactsPage;
