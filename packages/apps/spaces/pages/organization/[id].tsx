import React, { useEffect } from 'react';
import { DetailsPageLayout } from '@spaces/layouts/details-page-layout';
import styles from './organization.module.scss';
import { NextPageContext } from 'next';
import {
  ApolloClient,
  from,
  gql,
  HttpLink,
  InMemoryCache,
} from '@apollo/client';
import { authLink } from '../../apollo-client';
import { useRecoilState, useSetRecoilState } from 'recoil';
import Head from 'next/head';
import dynamic from 'next/dynamic';
import { showLegacyEditor } from '../../state/editor';
import { OrganizationDetailsSkeleton } from '@spaces/organization/organization-details/skeletons';
import { NoteEditorModes } from '@spaces/organization/editor/types';
import { TimelineSkeleton } from '@spaces/organisms/timeline/skeletons/TimelineSkeleton';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { organizationDetailsEdit } from '../../state';
import { TimelineContextProvider } from '@spaces/organisms/timeline/context/timelineContext';
import { IssueSummaryByStatus, ExternalSystem } from '@spaces/graphql';

const OrganizationTimeline = dynamic(
  () =>
    import('../../components/organization/organization-timeline').then(
      (res) => res.OrganizationTimeline,
    ),
  {
    ssr: true,
    loading: () => {
      return <TimelineSkeleton />;
    },
  },
);
const OrginizationToolbelt = dynamic(() =>
  import(
    '@spaces/organization/organization-toolbelt/OrginizationToolbelt'
  ).then((res) => res.OrginizationToolbelt),
);

const OrganizationEditor = dynamic(() =>
  import('@spaces/organization/editor/OrganizationEditor').then(
    (res) => res.OrganizationEditor,
  ),
);

const OrganizationDetails = dynamic(
  () =>
    import(
      '@spaces/organization/organization-details/OrganizationDetails'
    ).then((res) => res.OrganizationDetails),
  {
    loading: () => {
      return <OrganizationDetailsSkeleton />;
    },
  },
);

const OrganizationIssues = dynamic(() =>
  import('@spaces/organization/organization-issues/OrganizationIssues').then(
    (res) => res,
  ),
);

export async function getServerSideProps(context: NextPageContext) {
  const ssrClient = new ApolloClient({
    ssrMode: true,
    cache: new InMemoryCache(),
    link: from([
      authLink,
      new HttpLink({
        uri: `${process.env.SSR_PUBLIC_PATH}/customer-os-api/query`,
        fetchOptions: {
          credentials: 'include',
        },
      }),
    ]),
    queryDeduplication: true,
    assumeImmutableResults: true,
    connectToDevTools: true,
    credentials: 'include',
  });

  const organizationId = context.query.id;

  try {
    const res = await ssrClient.query({
      query: gql`
        query organization($id: ID!) {
          organization(id: $id) {
            id
            name
            issueSummaryByStatus {
              status
              count
            }
            externalLinks {
              type
              externalId
              externalSource
              externalUrl
            }
          }
        }
      `,
      variables: {
        id: organizationId,
      },
      context: {
        headers: {
          ...context?.req?.headers,
        },
      },
    });

    return {
      props: {
        name: res.data.organization.name || '',
        id: organizationId,
        issuSummaryByStatus: res.data.organization.issueSummaryByStatus,
        externalLinks: res.data.organization.externalLinks,
      },
    };
  } catch (e) {
    return {
      notFound: true,
    };
  }
}

function OrganizationDetailsPage({
  id,
  name,
  issuSummaryByStatus,
  externalLinks,
}: {
  id: string;
  name: string;
  issuSummaryByStatus: IssueSummaryByStatus[];
  externalLinks: ExternalSystem[];
}) {
  const [showEditor, setShowLegacyEditor] = useRecoilState(showLegacyEditor);
  const setOrganizationDetailsEdit = useSetRecoilState(organizationDetailsEdit);
  useEffect(() => {
    return () => {
      setShowLegacyEditor(false);
      setOrganizationDetailsEdit({ isEditMode: false });
    };
  }, []);

  return (
    <>
      <Head>
        <title>{!name ? 'Unnamed' : name}</title>
      </Head>
      <PageContentLayout>
        <DetailsPageLayout>
          <OrganizationDetails id={id}>
            <OrganizationIssues
              issueSummary={issuSummaryByStatus}
              externalLinks={externalLinks}
            />
          </OrganizationDetails>

          <TimelineContextProvider>
            {/*<section className={styles.notes}>*/}
            {/*  {!showEditor && <OrginizationToolbelt organizationId={id} />}*/}
            {/*  {showEditor && (*/}
            {/*    <OrganizationEditor*/}
            {/*      organizationId={id}*/}
            {/*      mode={NoteEditorModes.ADD}*/}
            {/*    />*/}
            {/*  )}*/}
            {/*</section>*/}
            <section className={styles.timeline}>
              <OrganizationTimeline id={id} />
            </section>
          </TimelineContextProvider>
        </DetailsPageLayout>
      </PageContentLayout>
    </>
  );
}

export default OrganizationDetailsPage;
