import React, { useEffect } from 'react';
import { DetailsPageLayout } from '../../components';
import styles from './contact.module.scss';
import { useRouter } from 'next/router';
import {
  ContactCommunicationDetails,
  ContactDetails,
  ContactEditor,
} from '../../components/contact';
import ContactHistory from '../../components/contact/contact-history/ContactHistory';
import { useSetRecoilState } from 'recoil';
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
import { ContactToolbelt } from '../../components/contact/contact-toolbelt/ContactToolbelt';
import { getContactPageTitle } from '../../utils';
import { Contact } from '../../graphQL/__generated__/generated';

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
  if (contactId == 'new') {
    // mutation
    const {
      data: { contact_Create },
    } = await ssrClient.mutate({
      mutation: gql`
        mutation createContact {
          contact_Create(input: { firstName: "", lastName: "" }) {
            id
            firstName
            lastName
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
        destination: `/contact/${contact_Create?.id}`,
      },
      props: {
        isEditMode: true,
        id: contact_Create?.id,
      },
    };
  }

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
  const setContactDetailsEdit = useSetRecoilState(contactDetailsEdit);
  useEffect(() => {
    setContactDetailsEdit({ isEditMode });
  }, [id, isEditMode]);

  return (
    <>
      <Head>
        <title> {getContactPageTitle(contact)}</title>
      </Head>
      <DetailsPageLayout onNavigateBack={() => push('/')}>
        <section className={styles.details}>
          <ContactDetails id={id as string} />
          <ContactCommunicationDetails id={id as string} />
        </section>
        <section className={styles.timeline}>
          <ContactHistory id={id as string} />
        </section>
        <section className={styles.notes}>
          <ContactEditor contactId={id as string} />
          <ContactToolbelt contactId={id} />
        </section>
      </DetailsPageLayout>
    </>
  );
}

export default ContactDetailsPage;
