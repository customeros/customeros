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
import { useRecoilState, useSetRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../state';
import Head from 'next/head';
import dynamic from 'next/dynamic';
import { NoteEditorModes } from '@spaces/organization/editor/OrganizationEditor';
import { OrganizationDetails } from '@spaces/organization/organization-details/OrganizationDetails';
import { OrginizationToolbelt } from '@spaces/organization/organization-toolbelt/OrginizationToolbelt';
import { showLegacyEditor } from '../../state/editor';
import { useAutoAnimate } from '@formkit/auto-animate/react';

// TODO add skeleton loader in options
const OrganizationContacts = dynamic(
  () =>
    import('../../components/organization').then(
      (res) => res.OrganizationContacts,
    ),
  { ssr: true },
);

const OrganizationTimeline = dynamic(
  () =>
    import('../../components/organization/organization-timeline').then(
      (res) => res.OrganizationTimeline,
    ),
  { ssr: true },
);
const OrganizationEditor = dynamic(() =>
  import('@spaces/organization/editor/OrganizationEditor').then(
    (res) => res.OrganizationEditor,
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
  if (organizationId == 'new') {
    // mutation
    const {
      data: { organization_Create },
    } = await ssrClient.mutate({
      mutation: gql`
        mutation createOrganization {
          organization_Create(input: { name: "" }) {
            id
          }
        }
      `,
      context: {
        headers: {
          ...context?.req?.headers,
        },
      },
    });

    return {
      redirect: {
        permanent: false,
        destination: `organization/${organization_Create?.id}`,
      },
      props: {
        isEditMode: true,
        id: organization_Create?.id,
      },
    };
  }

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
        isEditMode: !res.data.organization?.name.length,
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
  isEditMode,
  name,
}: {
  id: string;
  isEditMode: boolean;
  name: string;
}) {
  const { push } = useRouter();
  const setContactDetailsEdit = useSetRecoilState(organizationDetailsEdit);
  const [showEditor, setShowLegacyEditor] = useRecoilState(showLegacyEditor);
  const [animateRef] = useAutoAnimate({
    easing: 'ease-in',
  });
  useEffect(() => {
    setContactDetailsEdit({ isEditMode });
  }, [id, isEditMode]);

  useEffect(() => {
    return () => {
      setShowLegacyEditor(false);
    };
  }, []);

  return (
    <>
      <Head>
        <title>{isEditMode ? 'Unnamed' : name}</title>
      </Head>
      <DetailsPageLayout onNavigateBack={() => push('/')}>
        <section className={styles.organizationIdCard}>
          <OrganizationDetails id={id as string} />
        </section>
        <section className={styles.organizationDetails}>
          <OrganizationContacts id={id as string} />
        </section>
        <section className={styles.notes} ref={animateRef}>
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
