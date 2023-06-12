import type { NextPage } from 'next';
import React from 'react';
import Head from 'next/head';

import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { Filter } from '@spaces/graphql';

const CustomersPage: NextPage = () => {
  const preFilters = [
    {
      filter: {
        property: 'RELATIONSHIP',
        operation: 'EQ',
        value: ['CUSTOMER'],
      } as Filter,
    } as Filter,
  ];

  return (
    <>
      <Head>
        <title>Customers</title>
      </Head>
      <PageContentLayout>
        <OrganizationList preFilters={preFilters} />
      </PageContentLayout>
    </>
  );
};

export default CustomersPage;
