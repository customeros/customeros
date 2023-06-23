import React from 'react';
import { Skeleton } from '@spaces/atoms/skeleton';
import dynamic from 'next/dynamic';
import { Loader } from '@spaces/atoms/loader';

const OrganizationProfileSkeleton = dynamic(
  () =>
    import('@spaces/organization/skeletons/OrganizationProfileSkeleton').then(
      (res) => res.OrganizationProfileSkeleton,
    ),
  { ssr: false, loading: () => <Loader /> },
);
const ContactProfileSkeleton = dynamic(
  () =>
    import('@spaces/contact/skeletons/ContactProfileSkeleton').then(
      (res) => res.ContactProfileSkeleton,
    ),
  {
    ssr: false,
    loading: () => <Loader />,
  },
);

const TableSkeleton = dynamic(
  () =>
    import('@spaces/atoms/table/skeletons/TableSkeleton').then(
      (res) => res.TableSkeleton,
    ),
  {
    ssr: false,
    loading: () => <Loader />,
  },
);

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
  if (loadingUrl.match('contact') !== null) {
    return (
      <article style={{ gridArea: 'content' }}>
        <div style={{ margin: 'var(--spacing-sm) 0' }}>
          <Skeleton height='40px' width='60%' isSquare />
        </div>

        <TableSkeleton columns={4} />
      </article>
    );
  }
  if (
    loadingUrl.match('customers') !== null ||
    loadingUrl.match('organization')
  ) {
    return (
      <article style={{ gridArea: 'content' }}>
        <div style={{ margin: 'var(--spacing-sm) 0' }}>
          <Skeleton height='40px' width='60%' isSquare />
        </div>

        <TableSkeleton columns={5} />
      </article>
    );
  }

  return <Loader />;
};
