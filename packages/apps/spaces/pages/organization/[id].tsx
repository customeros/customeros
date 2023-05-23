import React, { useEffect } from 'react';
import { DetailsPageLayout } from '@spaces/layouts/details-page-layout';
import styles from './organization.module.scss';
import { useRouter } from 'next/router';
import { NextPageContext } from 'next';
import {
  ApolloClient,
  from,
  gql,
  HttpLink,
  InMemoryCache,
} from '@apollo/client';
import { authLink } from '../../apollo-client';
import { useRecoilState } from 'recoil';
import Head from 'next/head';
import dynamic from 'next/dynamic';
import { showLegacyEditor } from '../../state/editor';
import { OrganizationDetailsSkeleton } from '@spaces/organization/organization-details/skeletons';
import { NoteEditorModes } from '@spaces/organization/editor/types';
import { OrganizationContactsSkeleton } from '@spaces/organization/organization-contacts/skeletons';
import { TimelineSkeleton } from '@spaces/organisms/timeline';

// TODO add skeleton loader in options
const OrganizationContacts = dynamic(
  () =>
    import('../../components/organization').then(
      (res) => res.OrganizationContacts,
    ),
  {
    ssr: true,
    loading: () => <OrganizationContactsSkeleton />,
  },
);

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
}: {
  id: string;
  name: string;
}) {
  const { push } = useRouter();
  const [showEditor, setShowLegacyEditor] = useRecoilState(showLegacyEditor);


  useEffect(() => {
    return () => {
      setShowLegacyEditor(false);
    };
  }, []);

  return (
    <>
      <Head>
        <title>{!name ? 'Unnamed' : name}</title>
      </Head>
      <DetailsPageLayout>
        <section className={styles.organizationIdCard}>
          <OrganizationDetails id={id as string} />
        </section>
        <section className={styles.organizationDetails}>
          <OrganizationContacts id={id as string} />
        </section>
        <section className={styles.notes}>
          {!showEditor && <OrginizationToolbelt organizationId={id} />}
          {showEditor && (
            <OrganizationEditor
              organizationId={id as string}
              mode={NoteEditorModes.ADD}
            />
          )}
        </section>
        <section className={styles.timeline}>
          <OrganizationTimeline id={id as string} />
        </section>
      </DetailsPageLayout>
    </>
  );
}

export default OrganizationDetailsPage;
