import { useState, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { LoadingScreen } from '@shared/components/SplashScreen/components';
import { TimelineContextsProvider } from '@organization/components/TimelineContextsProvider';

import { SideSection } from './src/components/SideSection';
import { MainSection } from './src/components/MainSection';
import { Panels, TabsContainer } from './src/components/Tabs';
import { OrganizationTimelineWithActionsContext } from './src/components/Timeline/OrganizationTimelineWithActionsContext';

export const OrganizationPage = observer(() => {
  const navigate = useNavigate();
  const params = useParams();
  const [searchParams] = useSearchParams();
  const [isLoading, setIsLoading] = useState(true);

  const store = useStore();

  const { id } = params;

  if (typeof id === 'undefined') {
    navigate('/finder');

    return;
  }

  useEffect(() => {
    if (
      !store.organizations.value.has(id) &&
      !store.organizations.isFullyLoaded
    ) {
      setIsLoading(true);
      store.organizations.getById(id).then(() => {
        setIsLoading(false);
      });

      return;
    }
    setIsLoading(false);
  }, []);

  if (isLoading) {
    return <LoadingScreen hide={false} isLoaded={false} showSplash={true} />;
  }

  if (!isLoading && !store.organizations.value.has(id)) {
    throw new Error('Organization not found');
  }

  return (
    <div className='flex h-full'>
      <TimelineContextsProvider id={id}>
        <SideSection>
          <TabsContainer>
            <Panels tab={searchParams.get('tab') ?? 'about'} />
          </TabsContainer>
        </SideSection>

        <MainSection>
          <OrganizationTimelineWithActionsContext />
        </MainSection>
      </TimelineContextsProvider>
    </div>
  );
});
