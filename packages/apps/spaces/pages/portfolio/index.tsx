import type { NextPage } from 'next';
import React from 'react';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { useRecoilState, useRecoilValue } from 'recoil';
import { userData } from '../../state';
import { Filter } from '@spaces/graphql';
import { globalCacheData } from '../../state/globalCache';

const MyPortfolioPage: NextPage = () => {
  const { userId } = useRecoilValue(globalCacheData);
  const preFilters = [
    {
      filter: {
        property: 'OWNER_ID',
        operation: 'EQ',
        value: userId,
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
