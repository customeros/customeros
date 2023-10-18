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

  if (!issues.length) {
    return (
      <OrganizationPanel title='Issues' withFade>
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
    <OrganizationPanel title='Issues' withFade>
      <Flex as='article' w='full' direction='column'>
        <Heading fontWeight='semibold' fontSize='md' mb={2}>
          Open
        </Heading>
        <VStack>
          {!!openIssues?.length &&
            openIssues.map((issue, index) => (
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
      {!openIssues.length && (
        <EmptyIssueMessage
          description={`It looks like ${
            data?.organization?.name ?? '[Unknown]'
          } has no open issues at the moment`}
        />
      )}
      {!!closedIssues.length && (
        <Flex as='article' w='full' direction='column' mt={2}>
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
              {!closedIssues.length && (
                <EmptyIssueMessage
                  description={`It looks like ${
                    data?.organization?.name ?? '[Unknown]'
                  } has no closed issues at the moment`}
                />
              )}
              {!!closedIssues?.length && (
                <VStack>
                  {closedIssues.map((issue, index) => (
                    <IssueCard
                      issue={issue as any}
                      key={`issue-panel-${issue.id}`}
                    />
                  ))}
                </VStack>
              )}
            </Fade>
          </Collapse>
        </Flex>
      )}
    </OrganizationPanel>
  );
};
