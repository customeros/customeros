import React, { useEffect } from 'react';
import { DetailsPageLayout } from '../../components';
import styles from './organization.module.scss';
import { useRouter } from 'next/router';
import {
  OrganizationDetails,
  OrganizationEditor,
  NoteEditorModes,
  OrganizationContacts,
} from '../../components/organization';
import { OrganizationTimeline } from '../../components/organization/organization-timeline';
import { NextPageContext } from 'next';
import {
  ApolloClient,
  from,
  gql,
  HttpLink,
  InMemoryCache,
} from '@apollo/client';
import { authLink } from '../../apollo-client';
import { useSetRecoilState } from 'recoil';
import { contactDetailsEdit, organizationDetailsEdit } from '../../state';

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

    console.log('🏷️ ----- organization_Create: ', organization_Create);

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
        isEditMode:
          !res.data.organization?.name.length &&
          !res.data.organization?.name.length,
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
}: {
  id: string;
  isEditMode: boolean;
}) {
  const { push } = useRouter();
  const setContactDetailsEdit = useSetRecoilState(organizationDetailsEdit);

  useEffect(() => {
    if (isEditMode) {
      setContactDetailsEdit({ isEditMode: true });
    }
  }, [id, isEditMode]);

  return (
    <DetailsPageLayout onNavigateBack={() => push('/')}>
      <section className={styles.organizationIdCard}>
        <OrganizationDetails id={id as string} />
      </section>
      <section className={styles.organizationDetails}>
        <OrganizationContacts id={id as string} />
      </section>
      <section className={styles.notes}>
        <OrganizationEditor
          organizationId={id as string}
          mode={NoteEditorModes.ADD}
        />
      </section>
      <section className={styles.timeline}>
        <OrganizationTimeline id={id as string} />
      </section>
    </DetailsPageLayout>
  );
}

export default OrganizationDetailsPage;
