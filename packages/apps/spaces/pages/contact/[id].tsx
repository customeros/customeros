import React, { useEffect } from 'react';
import { DetailsPageLayout } from '@spaces/layouts/details-page-layout';
import styles from './contact.module.scss';
import { useRecoilState } from 'recoil';
import { authLink } from '../../apollo-client';
import {
  ApolloClient,
  from,
  gql,
  HttpLink,
  InMemoryCache,
} from '@apollo/client';
import { NextPageContext } from 'next';
import Head from 'next/head';
import { getContactPageTitle } from '../../utils/getContactPageTitle';
import { Contact } from '../../graphQL/__generated__/generated';
import { showLegacyEditor } from '../../state/editor';
import dynamic from 'next/dynamic';
import { ContactToolbelt } from '@spaces/contact/contact-toolbelt/ContactToolbelt';
import { ContactDetails } from '@spaces/contact/contact-details/ContactDetails';
import { ContactCommunicationDetails } from '@spaces/contact/contact-communication-details/ContactCommunicationDetails';
import { ContactLocations } from '@spaces/contact/contact-locations';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';

const ContactHistory = dynamic(
  () =>
    import('@spaces/contact/contact-history/ContactHistory').then(
      (res) => res.ContactHistory,
    ),
  { ssr: false },
);

const ContactEditor = dynamic(
  () =>
    import('@spaces/contact/editor/ContactEditor').then(
      (res) => res.ContactEditor,
    ),
  { ssr: false },
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
  const contactId = context.query.id;

  try {
    const res = await ssrClient.query({
      query: gql`
        query contact($id: ID!) {
          contact(id: $id) {
            id
            firstName
            lastName
            name
            emails {
              email
            }
            phoneNumbers {
              rawPhoneNumber
              e164
            }
            jobRoles {
              organization {
                name
              }
            }
          }
        }
      `,
      variables: {
        id: contactId,
      },
      context: {
        headers: {
          ...context?.req?.headers,
        },
      },
    });

    const contact = res.data?.contact;

    return {
      props: {
        id: contactId,
        contact,
      },
    };
  } catch (e) {
    return {
      notFound: true,
    };
  }
}
function ContactDetailsPage({ id, contact }: { id: string; contact: Contact }) {
  const [showEditor, setShowLegacyEditor] = useRecoilState(showLegacyEditor);

  useEffect(() => {
    return () => {
      setShowLegacyEditor(false);
    };
  }, []);

  return (
    <>
      <Head>
        <title> {getContactPageTitle(contact)}</title>
      </Head>
      <PageContentLayout>
        <DetailsPageLayout>
          <section className={styles.details}>
            <ContactDetails id={id} />
            <ContactCommunicationDetails id={id} />
            <ContactLocations id={id} />
          </section>
          <section className={styles.timeline}>
            <ContactHistory id={id} />
          </section>
          <section className={styles.notes}>
            {!showEditor && (
              <ContactToolbelt contactId={id} isSkewed={!showEditor} />
            )}
            {showEditor && <ContactEditor contactId={id} />}
          </section>
        </DetailsPageLayout>
      </PageContentLayout>
    </>
  );
}

export default ContactDetailsPage;
