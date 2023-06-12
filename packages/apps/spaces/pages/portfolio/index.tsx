import type { NextPage } from 'next';
import React from 'react';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { useRecoilState, useRecoilValue } from 'recoil';
import { userData } from '../../state';
import { Filter } from '@spaces/graphql';

const MyPortfolioPage: NextPage = () => {
  const loggedInUser = useRecoilValue(userData);
  const preFilters = [
    {
      filter: {
        property: 'OWNER_ID',
        operation: 'EQ',
        value: loggedInUser.identity,
      } as Filter,
    } as Filter,
  ];
  return (
    <>
      <Head>
        <title>My portfolio</title>
      </Head>
      <PageContentLayout>
        <OrganizationList preFilters={preFilters} />
      </PageContentLayout>
    </>
  );
};

export default MyPortfolioPage;
