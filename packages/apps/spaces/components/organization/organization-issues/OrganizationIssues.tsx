import {
  IssueSummaryByStatus,
  ExternalSystem,
  ExternalSystemType,
} from '@spaces/graphql';
import Link from 'next/link';
import { getZendeskBaseUrl } from '@spaces/utils/getZendeskBaseUrl';

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
    issueSummary.find((item) => item.status !== 'solved')?.count ?? 0;
  const zendesk = externalLinks.find(
    (url) => url.type === ExternalSystemType.ZendeskSupport,
  );
  const zendeskIssueId = zendesk?.externalId ?? '';
  const zendeskApiUrl = zendesk?.externalUrl ?? 'https://www.zendesk.com';

  const zendeskUrl = `${getZendeskBaseUrl(
    zendeskApiUrl,
  )}/${zendeskIssueId}/requester/requested_tickets`;
  const issueLabel = openIssuesCount === 1 ? 'issue' : 'issues';

  if (!openIssuesCount) {
    return null;
  }

  return (
    <article>
      <h1 className={styles.issuesHeader}>Issues</h1>
      <p className={styles.issuesItem}>
        {openIssuesCount} open {issueLabel}{' '}
        <Link href={zendeskUrl} target='_blank'>
          in Zendesk
        </Link>
      </p>
    </article>
  );
};

export default OrganizationIntegrations;
