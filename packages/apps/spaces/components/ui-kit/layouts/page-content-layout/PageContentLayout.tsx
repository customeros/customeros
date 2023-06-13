import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode, useEffect, useLayoutEffect } from 'react';
import styles from './page-content-layout.module.scss';
import { SidePanel } from '@spaces/organisms/side-panel';
import {useRecoilState, useSetRecoilState} from 'recoil';
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
  const [globalCache] = useRecoilState(globalCacheData);
  const setGlobalCacheData = useSetRecoilState(globalCacheData);
  const { onLoadGlobalCache, loading } = useGlobalCache();

  useEffect(() => {
    if (!globalCache.userId && !loading) {
      onLoadGlobalCache().then((res) => {
        setGlobalCacheData(res.data.global_Cache);
      });
    }
  }, [globalCache]);

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
