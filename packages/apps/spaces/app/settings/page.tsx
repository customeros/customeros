import React from 'react';

import {TabsContainer} from "./Tabs/TabsContainer";
import {Panels} from "./Tabs/Panels";
import {SettingsMainSection} from "./SettingsMainSection";

interface SettingsPageProps {
    searchParams: { tab?: string };
}

export default async function SettingsPage({searchParams}: SettingsPageProps) {

    return (
        <>
            <SettingsMainSection>
                <TabsContainer>
                    <Panels tab={searchParams.tab ?? 'oauth'} />
                </TabsContainer>
            </SettingsMainSection>
        </>
    );
};
