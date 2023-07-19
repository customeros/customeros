import React, { FC, ReactNode, useEffect } from 'react';
import { useRecoilState, useSetRecoilState } from 'recoil';
import { globalCacheData } from '../../../../state/globalCache';
import { useGlobalCache } from '@spaces/hooks/useGlobalCache';

import { PageLayout } from '@shared/components/PageLayout/PageLayout';

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

  return <PageLayout isOwner={globalCache?.isOwner}>{children}</PageLayout>;
};
