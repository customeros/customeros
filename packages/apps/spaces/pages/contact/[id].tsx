import React, { useEffect } from 'react';
import { DetailsPageLayout } from '@spaces/layouts/details-page-layout';
import styles from './contact.module.scss';
import { useRouter } from 'next/router';
import { useRecoilState, useSetRecoilState } from 'recoil';
import { contactDetailsEdit } from '../../state';
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
import { getContactPageTitle } from '../../utils';
import { Contact } from '../../graphQL/__generated__/generated';
import { showLegacyEditor } from '../../state/editor';
import { useAutoAnimate } from '@formkit/auto-animate/react';
import dynamic from 'next/dynamic';
import { ContactToolbelt } from '@spaces/contact/contact-toolbelt/ContactToolbelt';
import { ContactDetails } from '@spaces/contact/contact-details/ContactDetails';
import { ContactCommunicationDetails } from '@spaces/contact/contact-communication-details/ContactCommunicationDetails';

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
        isEditMode:
          !contact?.firstName?.length &&
          !contact?.lastName?.length &&
          !contact?.name?.length,
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
function ContactDetailsPage({
  id,
  isEditMode,
  contact,
}: {
  id: string;
  isEditMode: boolean;
  contact: Contact;
}) {
  const { push } = useRouter();
  const [showEditor, setShowLegacyEditor] = useRecoilState(showLegacyEditor);
  const [animateRef] = useAutoAnimate({
    easing: 'ease-in',
  });
  const setContactDetailsEdit = useSetRecoilState(contactDetailsEdit);
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
        <title> {getContactPageTitle(contact)}</title>
      </Head>
      <DetailsPageLayout onNavigateBack={() => push('/contact')}>
        <section className={styles.details}>
          <ContactDetails id={id as string} />
          <ContactCommunicationDetails id={id as string} />
        </section>
        <section className={styles.timeline}>
          <ContactHistory id={id as string} />
        </section>
        <section ref={animateRef} className={styles.notes}>
          {!showEditor && (
            <ContactToolbelt contactId={id} isSkewed={!showEditor} />
          )}
          {showEditor && <ContactEditor contactId={id} />}
        </section>
      </DetailsPageLayout>
    </>
  );
}

export default ContactDetailsPage;
