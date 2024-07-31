import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useConnections } from '@integration-app/react';

import { Issue } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { ChevronUp } from '@ui/media/icons/ChevronUp';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { IssueCard } from '@organization/components/Tabs/panels/IssuesPanel/IssueCard/IssueCard';
import { IssuesPanelSkeleton } from '@organization/components/Tabs/panels/IssuesPanel/IssuesPanelSkeleton';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';
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

export const IssuesPanel = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const organization = store.organizations.value.get(id);
  const issues = (
    store.timelineEvents.issues.getByOrganizationId(id) ?? []
  ).map((item) => item.value);

  const [isExpanded, setIsExpanded] = useState(true);
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

  if (loading && !store.demoMode) {
    return <IssuesPanelSkeleton />;
  }

  if (!connections.length) {
    return (
      <OrganizationPanel withFade title='Issues'>
        <EmptyIssueMessage title='Connect your customer support app'>
          To see your customers support issues here,{' '}
          <Link
            color='primary.600'
            className='text-primary-600'
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
        withFade
        title='Issues'
        actionItem={<ChannelLinkSelect />}
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
        withFade
        title='Issues'
        actionItem={<ChannelLinkSelect />}
      >
        <EmptyIssueMessage
          title='No issues detected'
          description={`It looks like ${
            organization?.value?.name ?? '[Unknown]'
          } has had a smooth journey thus far. Or
      perhaps they’ve been shy about reporting issues. Stay proactive and keep
      monitoring for optimal support.`}
        />
      </OrganizationPanel>
    );
  }

  return (
    <OrganizationPanel
      withFade
      title='Issues'
      actionItem={<ChannelLinkSelect />}
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
            organization?.value?.name ?? '[Unknown]'
          } has no open issues at the moment`}
        />
      )}
      {!!closedIssues.length && (
        <CollapsibleRoot
          open={isExpanded}
          className='flex flex-col w-full mt-2'
          onOpenChange={(value) => setIsExpanded(value)}
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
                  organization?.value?.name ?? '[Unknown]'
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
});
