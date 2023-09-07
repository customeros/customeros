import type { NextPage } from 'next';
import React from 'react';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';
import Head from 'next/head';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { useRecoilValue } from 'recoil';
import { Filter } from '@spaces/graphql';
import { globalCacheData } from '../../state/globalCache';
import { Portfolio } from '@spaces/atoms/icons';

const MyPortfolioPage: NextPage = () => {
  const { user } = useRecoilValue(globalCacheData);
  const preFilters = [
    {
      filter: {
        property: 'OWNER_ID',
        operation: 'EQ',
        value: user.id,
      } as Filter,
    } as Filter,
  ];
  const label =
    (user.firstName ? user.firstName + ' ' + user.lastName + "'s " : '') +
    'Portfolio';
  return (
    <>
      <Head>
        <title>My portfolio</title>
      </Head>
      <PageContentLayout>
        <OrganizationList
          icon={<Portfolio height={24} width={24} style={{ scale: '0.8' }} />}
          label={label}
          filterLabel={'portfolio'}
          preFilters={preFilters}
        />
      </PageContentLayout>
    </>
  );
};

export default MyPortfolioPage;
