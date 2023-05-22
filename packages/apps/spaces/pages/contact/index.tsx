import type { NextPage } from 'next';
import React from 'react';
import { PageContentLayout } from '../../components/ui-kit/layouts';
import Head from 'next/head';
import { ContactList } from '@spaces/contact/contact-list/ContactList';

const ContactsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Contacts</title>
      </Head>
      <PageContentLayout isSideBarShown={true}>
        <ContactList />
      </PageContentLayout>
    </>
  );
};

export default ContactsPage;
