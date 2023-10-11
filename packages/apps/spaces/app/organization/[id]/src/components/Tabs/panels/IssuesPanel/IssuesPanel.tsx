'use client';
import { useParams } from 'next/navigation';
import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Fade } from '@ui/transitions/Fade';
import { Heading } from '@ui/typography/Heading';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Issue } from '@graphql/types';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { IssueCard } from '@organization/src/components/Tabs/panels/IssuesPanel/IssueCard/IssueCard';
import { Collapse } from '@ui/transitions/Collapse';
import React, { useState } from 'react';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { ChevronUp } from '@ui/media/icons/ChevronUp';
import { IssuesPanelSkeleton } from '@organization/src/components/Tabs/panels/IssuesPanel/IssuesPanelSkeleton';
import { useGetIssuesQuery } from '@organization/src/graphql/getIssues.generated';
import { EmptyIssueMessage } from '@organization/src/components/Tabs/panels/IssuesPanel/EmptyIssueMessage/EmptyIssueMessage';

export const NEW_DATE = new Date(new Date().setDate(new Date().getDate() + 1));

// TODO uncomment commented out code as soon as COS-464 is merged

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
  const { data, isInitialLoading } = useGetIssuesQuery(client, {
    organizationId: id,
    from: NEW_DATE,
    size: 50,
  });
  const issues: Array<Issue> =
    (data?.organization?.timelineEvents as Array<Issue>) ?? [];
  const { open: openIssues, closed: closedIssues } = filterIssues(issues);

  if (isInitialLoading) {
    return <IssuesPanelSkeleton />;
  }
  return (
    <OrganizationPanel
      title='Issues'
      withFade
      bgImage={
        !issues?.length
          ? '/backgrounds/organization/half-circle-pattern.svg'
          : ''
      }
    >
      {!issues.length && (
        <EmptyIssueMessage
          organizationName={data?.organization?.name ?? '[Unknown]'}
        />
      )}

      {!!openIssues.length && (
        <Flex as='article' w='full' direction='column'>
          <Heading fontWeight='semibold' fontSize='md' mb={2}>
            Open
          </Heading>
          <VStack>
            {openIssues.map((issue, index) => (
              <Fade
                key={`issue-panel-${issue.id}`}
                in
                style={{ width: '100%' }}
              >
                <IssueCard issue={issue as any} />
              </Fade>
            ))}
          </VStack>
        </Flex>
      )}

      {!!closedIssues.length && (
        <Flex as='article' w='full' direction='column'>
          <Flex
            justifyContent='space-between'
            alignItems='center'
            w='full'
            as='button'
            pb={2}
            onClick={() => setIsExpanded((prev) => !prev)}
          >
            <Heading fontWeight='semibold' fontSize='md'>
              Closed
            </Heading>
            {isExpanded ? <ChevronDown /> : <ChevronUp />}
          </Flex>

          <Collapse
            in={isExpanded}
            style={{ overflow: 'unset' }}
            delay={{
              exit: 2,
            }}
          >
            <Fade
              in={isExpanded}
              delay={{
                enter: 0.2,
              }}
            >
              <VStack>
                {closedIssues.map((issue, index) => (
                  <IssueCard
                    issue={issue as any}
                    key={`issue-panel-${issue.id}`}
                  />
                ))}
              </VStack>
            </Fade>
          </Collapse>
        </Flex>
      )}
    </OrganizationPanel>
  );
};
