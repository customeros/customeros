import type { NextPage } from 'next';
import React from 'react';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { Company } from '@spaces/atoms/icons';

const OrganizationsPage: NextPage = () => {
  return (
    <>
      <Head>
        <title>Organizations</title>
      </Head>
      <PageContentLayout>
        <OrganizationList
          icon={<Company height={24} width={24} style={{ scale: '0.8' }} />}
          label={'Organizations'}
          filterLabel={'organizations'}
        />
      </PageContentLayout>
    </>
  );
};

export default OrganizationsPage;
