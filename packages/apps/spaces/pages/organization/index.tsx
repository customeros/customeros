import type { NextPage } from 'next';
import React from 'react';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';

const OrganizationsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Organizations</title>
      </Head>
      <PageContentLayout isSideBarShown={true}>
        <OrganizationList />
      </PageContentLayout>
    </>
  );
};

export default OrganizationsPage;
