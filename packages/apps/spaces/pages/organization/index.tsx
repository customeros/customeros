import type { NextPage } from 'next';
import React, { useState } from 'react';
import { PageContentLayout } from '../../components/ui-kit/layouts';
import { SidePanel } from '@spaces/organisms/side-panel/SidePanel';
import { OrganizationList } from '@spaces/organization/organization-list/OrganizationList';

const OrganizationsPage: NextPage = () => {
  const [isSidePanelVisible, setSidePanelVisible] = useState(false);

  return (
    <PageContentLayout isPanelOpen={isSidePanelVisible} isSideBarShown={true}>
      <SidePanel
        onPanelToggle={setSidePanelVisible}
        isPanelOpen={isSidePanelVisible}
      />
      <article style={{ gridArea: 'content' }}>
        <OrganizationList />
      </article>
    </PageContentLayout>
  );
};

export default OrganizationsPage;
