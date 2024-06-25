import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { TimelineContextsProvider } from '@organization/components/TimelineContextsProvider';

import { SideSection } from './src/components/SideSection';
import { MainSection } from './src/components/MainSection';
import { Panels, TabsContainer } from './src/components/Tabs';
import { OrganizationTimelineWithActionsContext } from './src/components/Timeline/OrganizationTimelineWithActionsContext';

export const OrganizationPage = () => {
  const navigate = useNavigate();
  const params = useParams();
  const [searchParams] = useSearchParams();

  const { id } = params;

  if (typeof id === 'undefined') {
    navigate('/finder');

    return;
  }

  // add logic to redirect if organization is hidden
  // if (organizationData?.hide) {
  //   notFound();
  // }

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
};
