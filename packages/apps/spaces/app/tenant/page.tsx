import type { NextPage } from 'next';
import React from 'react';

import {TabsContainer} from "./Tabs/TabsContainer";
import {Panels} from "./Tabs/Panels";
import {SideSection} from "./SideSection";
import {MainSection} from "./MainSection";

interface TenantPageProps {
    searchParams: { tab?: string };
}

export default async function TenantPage({searchParams}: TenantPageProps) {

    return (
        <>
            <MainSection>
                <TabsContainer>
                    <Panels tab={searchParams.tab ?? 'oauth'} />
                </TabsContainer>
            </MainSection>
        </>
    );
};
