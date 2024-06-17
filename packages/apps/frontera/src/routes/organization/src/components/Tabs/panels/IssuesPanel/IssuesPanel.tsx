import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';

import { useConnections } from '@integration-app/react';

import { Issue } from '@graphql/types';
import { ChevronUp } from '@ui/media/icons/ChevronUp';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetIssuesQuery } from '@organization/graphql/getIssues.generated';
import { IssueCard } from '@organization/components/Tabs/panels/IssuesPanel/IssueCard/IssueCard';
import { IssuesPanelSkeleton } from '@organization/components/Tabs/panels/IssuesPanel/IssuesPanelSkeleton';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import {
  CollapsibleRoot,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@ui/transitions/Collapse/Collapse';
import { EmptyIssueMessage } from '@organization/components/Tabs/panels/IssuesPanel/EmptyIssueMessage/EmptyIssueMessage';

import { ChannelLinkSelect } from './ChannelLinkSelect';

export const NEW_DATE = new Date(new Date().setDate(new Date().getDate() + 1));

function filterIssues(issues: Array<Issue>): {
  open: Array<Issue>;
  closed: Array<Issue>;
} {
  return issues.reduce(
    (
      acc: {
        open: Array<Issue>;
        closed: Array<Issue>;
      },
      issue,
    ) => {
      if (['closed', 'solved'].includes(issue.status.toLowerCase())) {
        acc.closed.push(issue);
      } else {
        acc.open.push(issue);
      }

      return acc;
    },
    { open: [], closed: [] },
  );
}

export const IssuesPanel = () => {
  const id = useParams()?.id as string;
  const client = getGraphQLClient();
  const [isExpanded, setIsExpanded] = useState(true);
  const { data, isLoading } = useGetIssuesQuery(client, {
    organizationId: id,
    from: NEW_DATE,
    size: 50,
  });
  const issues: Array<Issue> =
    (data?.organization?.timelineEvents as Array<Issue>) ?? [];
  const { open: openIssues, closed: closedIssues } = filterIssues(issues);
  const { items, loading } = useConnections();
  const connections = items
    .map((item) => item.integration?.key)
    .filter((item) =>
      [
        'unthread',
        'zendesk',
        'dixa',
        'slack',
        'pylon',
        'intercom',
        'crisp',
        'atlassian',
      ].includes(item ?? ''),
    );

  if (loading || isLoading) {
    return <IssuesPanelSkeleton />;
  }

  if (!connections.length) {
    return (
      <OrganizationPanel title='Issues' withFade>
        <EmptyIssueMessage title='Connect your customer support app'>
          To see your customers support issues here,{' '}
          <Link
            className='text-primary-600'
            color='primary.600'
            to='/settings?tab=integrations'
          >
            Go to settings
          </Link>{' '}
          and connect an app like Zendesk or Unthread.
        </EmptyIssueMessage>
      </OrganizationPanel>
    );
  }

  if (connections?.[0] === 'unthread' && !issues.length) {
    return (
      <OrganizationPanel
        title='Issues'
        withFade
        actionItem={<ChannelLinkSelect from={NEW_DATE} />}
      >
        <EmptyIssueMessage title='Link an Unthread Slack channel'>
          Show your Unthread support issues here by linking a Slack channel.
        </EmptyIssueMessage>
      </OrganizationPanel>
    );
  }

  if (!issues.length) {
    return (
      <OrganizationPanel
        title='Issues'
        withFade
        actionItem={<ChannelLinkSelect from={NEW_DATE} />}
      >
        <EmptyIssueMessage
          title='No issues detected'
          description={`It looks like ${
            data?.organization?.name ?? '[Unknown]'
          } has had a smooth journey thus far. Or
      perhaps they’ve been shy about reporting issues. Stay proactive and keep
      monitoring for optimal support.`}
        />
      </OrganizationPanel>
    );
  }

  return (
    <OrganizationPanel
      title='Issues'
      withFade
      actionItem={<ChannelLinkSelect from={NEW_DATE} />}
    >
      <article className='w-full flex flex-col'>
        <h2 className='text-base font-semibold mb-2'>Open</h2>
        <div className='flex flex-col gap-2'>
          {!!openIssues?.length &&
            openIssues.map((issue, index) => (
              <IssueCard key={index} issue={issue} />
            ))}
        </div>
      </article>
      {!openIssues.length && (
        <EmptyIssueMessage
          description={`It looks like ${
            data?.organization?.name ?? '[Unknown]'
          } has no open issues at the moment`}
        />
      )}
      {!!closedIssues.length && (
        <CollapsibleRoot
          open={isExpanded}
          onOpenChange={(value) => setIsExpanded(value)}
          className='flex flex-col w-full mt-2'
        >
          {isExpanded}
          <div className='flex justify-between w-full items-center pb-2'>
            <h2 className='font-semibold text-base'>Closed</h2>
            <CollapsibleTrigger asChild={false}>
              {isExpanded ? <ChevronDown /> : <ChevronUp />}
            </CollapsibleTrigger>
          </div>
          <CollapsibleContent>
            {!closedIssues.length && (
              <EmptyIssueMessage
                description={`It looks like ${
                  data?.organization?.name ?? '[Unknown]'
                } has no closed issues at the moment`}
              />
            )}
            {!!closedIssues?.length && (
              <div className='flex flex-col space-y-2'>
                {closedIssues.map((issue) => (
                  <IssueCard issue={issue} key={`issue-panel-${issue.id}`} />
                ))}
              </div>
            )}
          </CollapsibleContent>
        </CollapsibleRoot>
      )}
    </OrganizationPanel>
  );
};
