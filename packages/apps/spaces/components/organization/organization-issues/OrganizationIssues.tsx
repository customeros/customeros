import {
  IssueSummaryByStatus,
  ExternalSystem,
  ExternalSystemType,
} from '@spaces/graphql';
import Link from 'next/link';

import styles from './organization-issues.module.scss';

interface OrganizationIntegrationsProps {
  issueSummary: IssueSummaryByStatus[];
  externalLinks: ExternalSystem[];
}

const OrganizationIntegrations = ({
  issueSummary,
  externalLinks,
}: OrganizationIntegrationsProps) => {
  const openIssuesCount =
    issueSummary.find((item) => item.status === 'open')?.count ?? 0;
  const zendeskUrl =
    externalLinks.find((url) => url.type === ExternalSystemType.ZendeskSupport)
      ?.externalUrl ?? 'https://www.zendesk.com';

  if (!openIssuesCount) {
    return null;
  }

  return (
    <article>
      <h1 className={styles.issuesHeader}>Issues</h1>
      <p className={styles.issuesItem}>
        {openIssuesCount} open issues{' '}
        <Link href={zendeskUrl} target='_blank'>
          in Zendesk
        </Link>
      </p>
    </article>
  );
};

export default OrganizationIntegrations;
