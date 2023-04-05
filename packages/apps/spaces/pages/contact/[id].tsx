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
// import { authLink, httpLinkSSR } from '../../apollo-client';
import { ApolloClient, from, gql, InMemoryCache } from '@apollo/client';

export async function getServerSideProps(context) {
  const ssrClient = new ApolloClient({
    ssrMode: true,
    cache: new InMemoryCache(),
    // link: from([authLink, httpLinkSSR]),
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
          ...context.req.headers,
        },
      },
    });

    return {
      redirect: {
        permanent: false,
        destination: `http://localhost:3001/contact/${contact_Create?.id}`,
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
          }
        }
      `,
      variables: {
        id: contactId,
      },
      context: {
        headers: {
          ...context.req.headers,
        },
      },
    });
    return {
      props: {
        isEditMode:
          !res.data.contact.firstName.length &&
          !res.data.contact.firstName.length,
        id: contactId,
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
}: {
  id: string;
  isEditMode: boolean;
}) {
  const { push } = useRouter();
  const setContactDetailsEdit = useSetRecoilState(contactDetailsEdit);
  useEffect(() => {
    if (isEditMode) {
      setContactDetailsEdit({ isEditMode: true });
    }
  }, [id, isEditMode]);

  return (
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
      </section>
    </DetailsPageLayout>
  );
}

export default ContactDetailsPage;
