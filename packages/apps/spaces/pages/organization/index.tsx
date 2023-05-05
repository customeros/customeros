import type { NextPage } from 'next';
import React, { useState } from 'react';
import { PageContentLayout } from '../../components/ui-kit/layouts';
import { SidePanel } from '../../components/ui-kit/organisms';
import { OrganizationList } from '../../components/organization/organization-list/OrganizationList';

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
