import React, { FC, ReactNode, useEffect } from 'react';
import { useRecoilState, useSetRecoilState } from 'recoil';
import { globalCacheData } from '../../../../state/globalCache';
import { useGlobalCache } from '@spaces/hooks/useGlobalCache';

import { PageLayout } from 'app/components/PageLayout/PageLayout';
import { RootSidenav } from 'app/components/RootSidenav/RootSidenav';
import { GridItem } from '@ui/layout/Grid';
interface PageContentLayout {
  children: ReactNode;
}

export const PageContentLayout: FC<PageContentLayout> = ({ children }) => {
  const [globalCache] = useRecoilState(globalCacheData);
  const setGlobalCacheData = useSetRecoilState(globalCacheData);
  const { onLoadGlobalCache, loading } = useGlobalCache();

  useEffect(() => {
    if (!globalCache?.user?.id && !loading) {
      onLoadGlobalCache().then((res) => {
        setGlobalCacheData(res?.data?.global_Cache);
      });
    }
  }, [globalCache]);

  return (
    <PageLayout>
      <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
        {children}
      </GridItem>
    </PageLayout>
  );
};
