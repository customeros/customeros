import React from 'react';
import { OrganizationProfileSkeleton } from '@spaces/organization/skeletons/OrganizationProfileSkeleton';
import { ContactProfileSkeleton } from '@spaces/contact/skeletons/ContactProfileSkeleton';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { TableSkeleton } from '@spaces/atoms/table/skeletons/';
import { Skeleton } from '@spaces/atoms/skeleton';

interface PageSkeletonProps {
  loadingUrl: string;
}

export const PageSkeleton: React.FC<PageSkeletonProps> = ({ loadingUrl }) => {
  const organizationProfile =
    /^\/organization\/[0-9a-f]{8}-(?:[0-9a-f]{4}-){3}[0-9a-f]{12}$/i;
  const newOrganizationProfile = /^\/organization\/new/i;
  const contactProfile =
    /^\/contact\/[0-9a-f]{8}-(?:[0-9a-f]{4}-){3}[0-9a-f]{12}$/i;
  const contactNew = /^\/contact\/new/i;
  if (loadingUrl.match(organizationProfile) !== null) {
    return <OrganizationProfileSkeleton />;
  }
  if (loadingUrl.match(newOrganizationProfile) !== null) {
    return <OrganizationProfileSkeleton />;
  }
  if (loadingUrl.match(contactProfile) !== null) {
    return <ContactProfileSkeleton />;
  }
  if (loadingUrl.match(contactNew) !== null) {
    return <ContactProfileSkeleton />;
  }
  if (
    loadingUrl.match('contact') !== null ||
    loadingUrl.match('organization')
  ) {
    return (
      <PageContentLayout isSideBarShown={true}>
        <article style={{ gridArea: 'content' }}>
          <div style={{ margin: 'var(--spacing-sm) 0' }}>
            <Skeleton height='40px' width='60%' isSquare />
          </div>

          <TableSkeleton columns={5} />
        </article>
      </PageContentLayout>
    );
  }

  return (
    <div className='loader_container'>
      <div className='loader'>
        <div className='blue'/>
      </div>
    </div>
  );
};
