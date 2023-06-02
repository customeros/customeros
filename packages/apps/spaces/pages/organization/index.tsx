import type { NextPage } from 'next';
import React from 'react';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';

const OrganizationsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Organizations</title>
      </Head>
      <PageContentLayout>
        <OrganizationList />
      </PageContentLayout>
    </>
  );
};

export default OrganizationsPage;
