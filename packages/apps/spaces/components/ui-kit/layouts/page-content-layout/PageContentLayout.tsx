import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode, useLayoutEffect } from 'react';
import styles from './page-content-layout.module.scss';
import { SidePanel } from '@spaces/organisms/side-panel';
import { useSetRecoilState } from 'recoil';
import { globalCacheData } from '../../../../state/globalCache';
import { NextPageContext } from 'next';
import {
  ApolloClient,
  from,
  gql,
  HttpLink,
  InMemoryCache,
} from '@apollo/client';
import { authLink } from '../../../../apollo-client';
import { useGlobalCache } from '@spaces/hooks/useGlobalCache';

interface PageContentLayout {
  children: ReactNode;
}

export const PageContentLayout: FC<PageContentLayout> = ({ children }) => {
  // let setGlobalCacheData = useSetRecoilState(globalCacheData);
  // const { data, loading, error } = useGlobalCache();

  // useLayoutEffect(() => {
  //   if (!loading && data) {
  //     console.log('setting global cache data');
  //     console.log(data);
  //     setGlobalCacheData(data);
  //   }
  // }, [data, loading]);

  return (
    <div className={styles.pageContent}>
      <SidePanel />
      <div
        style={{
          padding: '1.2rem',
          height: '100%',
          gridArea: 'content',
          overflowX: 'hidden',
          overflowY: 'auto',
        }}
      >
        {children}
      </div>
    </div>
  );
};
