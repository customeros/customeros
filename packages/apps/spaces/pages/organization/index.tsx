import type { NextPage } from 'next';
import React from 'react';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';

const OrganizationsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Organizations</title>
      </Head>
      <OrganizationList />
    </>
  );
};

export default OrganizationsPage;
