import type {NextPage} from 'next';
import React, {useState} from 'react';
import {PageContentLayout} from '../../components/ui-kit/layouts';
import {SidePanel} from '../../components/ui-kit/organisms';
import {useRouter} from 'next/router';
import Head from 'next/head';
import {ContactList} from "../../components/contact/contact-list/ContactList";

const ContactsPage: NextPage = () => {
    const [isSidePanelVisible, setSidePanelVisible] = useState(false);

    return (
        <>
            <Head>
                <title>Contacts</title>
            </Head>
            <PageContentLayout isPanelOpen={isSidePanelVisible} isSideBarShown={true}>
                <SidePanel
                    onPanelToggle={setSidePanelVisible}
                    isPanelOpen={isSidePanelVisible}
                ></SidePanel>
                <article style={{gridArea: 'content'}}>
                    <ContactList/>
                </article>
            </PageContentLayout>
        </>
    );
};

export default ContactsPage;
