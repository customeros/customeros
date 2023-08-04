'use client';

import {Card, CardBody, CardHeader} from "@ui/layout/Card";
import React, {useEffect, useRef, useState} from "react";
import {Heading} from "@ui/typography/Heading";
import {SettingsIntegrationItem} from "@spaces/molecules/settings-integration-item";
import {GetSettings} from "../../../../../services";
import {toast} from "react-toastify";
import {DebouncedInput} from "@spaces/atoms/input";
import Search from "@spaces/atoms/icons/Search";
import {Skeleton} from "@spaces/atoms/skeleton";

export const IntegrationsPanel = () => {
    const [reload, setReload] = useState<boolean>(false);
    const reloadRef = useRef<boolean>(reload);

    const [loading, setLoading] = useState<boolean>(true);

    //states: ACTIVE OR INACTIVE
    //TODO: switch to a different state when the integration is being configured to fetch running and error states
    const [integrations, setIntegrations] = useState([
        {
            key: 'gsuite',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'/logos/google-icon.svg'}
                    identifier={'gsuite'}
                    name={'G Suite'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'privateKey',
                            label: 'Private key',
                            textarea: true,
                        },
                        {
                            name: 'clientEmail',
                            label: 'Service account email',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'hubspot',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'/logos/hubspot.svg'}
                    identifier={'hubspot'}
                    name={'Hubspot'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'privateAppKey',
                            label: 'API key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'smartsheet',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'/logos/smartsheet.svg'}
                    identifier={'smartsheet'}
                    name={'Smartsheet'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'id',
                            label: 'ID',
                        },
                        {
                            name: 'accessToken',
                            label: 'API key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'jira',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'/logos/jira.svg'}
                    identifier={'jira'}
                    name={'Jira'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                        {
                            name: 'email',
                            label: 'Email',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'trello',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'/logos/trello.svg'}
                    identifier={'trello'}
                    name={'Trello'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                        {
                            name: 'apiKey',
                            label: 'API key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'aha',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/aha.svg'}
                    identifier={'aha'}
                    name={'Aha'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiUrl',
                            label: 'API Url',
                        },
                        {
                            name: 'apiKey',
                            label: 'API key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'airtable',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/airtable.svg'}
                    identifier={'airtable'}
                    name={'Airtable'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'personalAccessToken',
                            label: 'Personal access token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'amplitude',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/amplitude.svg'}
                    identifier={'amplitude'}
                    name={'Amplitude'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API key',
                        },
                        {
                            name: 'secretKey',
                            label: 'Secret key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'asana',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/asana.svg'}
                    identifier={'asana'}
                    name={'Asana'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'baton',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'baton'}
                    name={'Baton'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'babelforce',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'babelforce'}
                    name={'Babelforce'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'regionEnvironment',
                            label: 'Region / Environment',
                        },
                        {
                            name: 'accessKeyId',
                            label: 'Access Key Id',
                        },
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'bigquery',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/bigquery.svg'}
                    identifier={'bigquery'}
                    name={'BigQuery'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'serviceAccountKey',
                            label: 'Service account key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'braintree',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/braintree.svg'}
                    identifier={'braintree'}
                    name={'Braintree'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'publicKey',
                            label: 'Public Key',
                        },
                        {
                            name: 'privateKey',
                            label: 'Private Key',
                        },
                        {
                            name: 'environment',
                            label: 'Environment',
                        },
                        {
                            name: 'merchantId',
                            label: 'Merchant Id',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'callrail',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'callrail'}
                    name={'CallRail'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'account',
                            label: 'Account',
                        },
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'chargebee',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/chargebee.svg'}
                    identifier={'chargebee'}
                    name={'Chargebee'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'productCatalog',
                            label: 'Product Catalog',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'chargify',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/chargify.svg'}
                    identifier={'chargify'}
                    name={'Chargify'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'clickup',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/clickup.svg'}
                    identifier={'clickup'}
                    name={'ClickUp'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'closecom',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/close.com.svg'}
                    identifier={'closecom'}
                    name={'Close.com'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'coda',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/coda.svg'}
                    identifier={'coda'}
                    name={'Coda'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'authToken',
                            label: 'Auth Token',
                        },
                        {
                            name: 'documentId',
                            label: 'Document Id',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'confluence',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/confluence.svg'}
                    identifier={'confluence'}
                    name={'Confluence'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                        {
                            name: 'loginEmail',
                            label: 'Login Email',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'courier',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/courier.svg'}
                    identifier={'courier'}
                    name={'Courier'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'customerio',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-io.svg'}
                    identifier={'customerio'}
                    name={'Customer.io'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'datadog',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/datadog.svg'}
                    identifier={'datadog'}
                    name={'Datadog'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'applicationKey',
                            label: 'Application Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'delighted',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/delighted.svg'}
                    identifier={'delighted'}
                    name={'Delighted'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'dixa',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/dixa.svg'}
                    identifier={'dixa'}
                    name={'Dixa'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'drift',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/drift.svg'}
                    identifier={'drift'}
                    name={'Drift'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'emailoctopus',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/emailoctopus.svg'}
                    identifier={'emailoctopus'}
                    name={'EmailOctopus'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'facebookMarketing',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/facebook.svg'}
                    identifier={'facebookMarketing'}
                    name={'Facebook Marketing'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'fastbill',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/fastbill.svg'}
                    identifier={'fastbill'}
                    name={'Fastbill'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'projectId',
                            label: 'Project Id',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'flexport',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/flexport.svg'}
                    identifier={'flexport'}
                    name={'Flexport'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'freshcaller',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/freshcaller.svg'}
                    identifier={'freshcaller'}
                    name={'Freshcaller'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'freshdesk',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/freshdesk.svg'}
                    identifier={'freshdesk'}
                    name={'Freshdesk'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'freshsales',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/freshsales.svg'}
                    identifier={'freshsales'}
                    name={'Freshsales'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'freshservice',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/freshservice.svg'}
                    identifier={'freshservice'}
                    name={'Freshservice'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'genesys',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'genesys'}
                    name={'Genesys'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'region',
                            label: 'Region',
                        },
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'github',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/github.svg'}
                    identifier={'github'}
                    name={'Github'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'gitlab',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/gitlab.svg'}
                    identifier={'gitlab'}
                    name={'GitLab'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'gocardless',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'gocardless'}
                    name={'GoCardless'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                        {
                            name: 'environment',
                            label: 'Environment',
                        },
                        {
                            name: 'version',
                            label: 'Version',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'gong',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'gong'}
                    name={'Gong'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'harvest',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/harvest.svg'}
                    identifier={'harvest'}
                    name={'Harvest'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accountId',
                            label: 'Account Id',
                        },
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'instagram',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/instagram.svg'}
                    identifier={'instagram'}
                    name={'Instagram'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'instatus',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'instatus'}
                    name={'Instatus'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'intercom',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/intercom.svg'}
                    identifier={'intercom'}
                    name={'Intercom'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'klaviyo',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'klaviyo'}
                    name={'Klaviyo'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'kustomer',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/kustomer.svg'}
                    identifier={'kustomer'}
                    name={'Kustomer'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'looker',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/looker.svg'}
                    identifier={'looker'}
                    name={'Looker'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                        {
                            name: 'domain',
                            label: 'Domain',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'mailchimp',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/mailchimp.svg'}
                    identifier={'mailchimp'}
                    name={'Mailchimp'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'mailjetemail',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/mailjetemail.svg'}
                    identifier={'mailjetemail'}
                    name={'Mailjet Email'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'apiSecret',
                            label: 'API Secret',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'marketo',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/marketo.svg'}
                    identifier={'marketo'}
                    name={'Marketo'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                        {
                            name: 'domainUrl',
                            label: 'Domain Url',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'microsoftteams',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/microsoftteams.svg'}
                    identifier={'microsoftteams'}
                    name={'Microsoft Teams'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'tenantId',
                            label: 'Tenant Id',
                        },
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'monday',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'monday'}
                    name={'Monday'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'paypaltransaction',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/paypal.svg'}
                    identifier={'paypaltransaction'}
                    name={'Paypal Transaction'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'secret',
                            label: 'Secret',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'notion',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/notion.svg'}
                    identifier={'notion'}
                    name={'Notion'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'internalAccessToken',
                            label: 'Internal Access Token',
                        },
                        {
                            name: 'publicClientId',
                            label: 'Public Client Id',
                        },
                        {
                            name: 'publicClientSecret',
                            label: 'Public Client Secret',
                        },
                        {
                            name: 'publicAccessToken',
                            label: 'Public Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'oraclenetsuite',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/oraclenetsuite.svg'}
                    identifier={'oraclenetsuite'}
                    name={'Oracle Netsuite'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accountId',
                            label: 'Account Id',
                        },
                        {
                            name: 'consumerKey',
                            label: 'Consumer Key',
                        },
                        {
                            name: 'consumerSecret',
                            label: 'Consumer Secret',
                        },
                        {
                            name: 'tokenId',
                            label: 'Token Id',
                        },
                        {
                            name: 'tokenSecret',
                            label: 'Token Secret',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'orb',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/orb.svg'}
                    identifier={'orb'}
                    name={'Orb'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'orbit',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/orbit.svg'}
                    identifier={'orbit'}
                    name={'Orbit'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'pagerduty',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/pagerduty.svg'}
                    identifier={'pagerduty'}
                    name={'PagerDuty'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'paystack',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'paystack'}
                    name={'Paystack'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'secretKey',
                            label: 'Secret Key (mandatory)',
                        },
                        {
                            name: 'lookbackWindow',
                            label: 'Lookback Window (in days)',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'pendo',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'pendo'}
                    name={'Pendo'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'pipedrive',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/pipedrive.svg'}
                    identifier={'pipedrive'}
                    name={'Pipedrive'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'plaid',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/plaid.svg'}
                    identifier={'plaid'}
                    name={'Plaid'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'plausible',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/plausible.svg'}
                    identifier={'plausible'}
                    name={'Plausible'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'siteId',
                            label: 'Site Id',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'posthog',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/posthog.svg'}
                    identifier={'posthog'}
                    name={'PostHog'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key (mandatory)',
                        },
                        {
                            name: 'baseUrl',
                            label: 'Base URL (optional)',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'qualaroo',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/qualaroo.svg'}
                    identifier={'qualaroo'}
                    name={'Qualaroo'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'quickbooks',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/quickbooks.svg'}
                    identifier={'quickbooks'}
                    name={'QuickBooks'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                        {
                            name: 'realmId',
                            label: 'Realm Id',
                        },
                        {
                            name: 'refreshToken',
                            label: 'Refresh Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'recharge',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/recharge.svg'}
                    identifier={'recharge'}
                    name={'Recharge'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'recruitee',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/recruitee.svg'}
                    identifier={'recruitee'}
                    name={'Recruitee'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'companyId',
                            label: 'Company Id',
                        },
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'recurly',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/recurly.svg'}
                    identifier={'recurly'}
                    name={'Recurly'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'retently',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/retently.svg'}
                    identifier={'retently'}
                    name={'Retently'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'salesforce',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/salesforce.svg'}
                    identifier={'salesforce'}
                    name={'Salesforce'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'clientId',
                            label: 'Client Id',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                        {
                            name: 'refreshToken',
                            label: 'Refresh Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'salesloft',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/salesloft.svg'}
                    identifier={'salesloft'}
                    name={'SalesLoft'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'sendgrid',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/sendgrid.svg'}
                    identifier={'sendgrid'}
                    name={'SendGrid'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'sentry',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/sentry.svg'}
                    identifier={'sentry'}
                    name={'Sentry'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'project',
                            label: 'Project',
                        },
                        {
                            name: 'authenticationToken',
                            label: 'Authentication Token',
                        },
                        {
                            name: 'organization',
                            label: 'Organization',
                        },
                        {
                            name: 'host',
                            label: 'Host (optional)',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'slack',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/slack.svg'}
                    identifier={'slack'}
                    name={'Slack'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                        {
                            name: 'channelFilter',
                            label: 'Channel Filter',
                        },
                        {
                            name: 'lookbackWindow',
                            label: 'lookback Window (in days, optional)',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'stripe',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/stripe.svg'}
                    identifier={'stripe'}
                    name={'Stripe'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accountId',
                            label: 'Account Id',
                        },
                        {
                            name: 'secretKey',
                            label: 'Secret Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'surveysparrow',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/surveysparrow.svg'}
                    identifier={'surveysparrow'}
                    name={'SurveySparrow'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'surveymonkey',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/surveymonkey.svg'}
                    identifier={'surveymonkey'}
                    name={'SurveyMonkey'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'talkdesk',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/talkdesk.svg'}
                    identifier={'talkdesk'}
                    name={'TalkDesk'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'tiktok',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/tiktok.svg'}
                    identifier={'tiktok'}
                    name={'TikTok'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'todoist',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'todoist'}
                    name={'Todoist'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'typeform',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'typeform'}
                    name={'Typeform'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'vittally',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'vittally'}
                    name={'Vittally'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'wrike',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/wrike.svg'}
                    identifier={'wrike'}
                    name={'Wrike'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'accessToken',
                            label: 'Access Token',
                        },
                        {
                            name: 'hostUrl',
                            label: 'Host URL',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'xero',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/customer-os-small.svg'}
                    identifier={'xero'}
                    name={'Xero'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'clientId',
                            label: 'Client ID',
                        },
                        {
                            name: 'clientSecret',
                            label: 'Client Secret',
                        },
                        {
                            name: 'tenantId',
                            label: 'Tenant ID',
                        },
                        {
                            name: 'scopes',
                            label: 'Scopes',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'zendesksupport',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'/logos/zendesksupport.svg'}
                    identifier={'zendesksupport'}
                    name={'Zendesk Support'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiKey',
                            label: 'API key',
                        },
                        {
                            name: 'subdomain',
                            label: 'Subdomain',
                        },
                        {
                            name: 'adminEmail',
                            label: 'Admin email',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'zendeskchat',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/zendeskchat.svg'}
                    identifier={'zendeskchat'}
                    name={'Zendesk Chat'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'subdomain',
                            label: 'Subdomain',
                        },
                        {
                            name: 'accessKey',
                            label: 'Access Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'zendesktalk',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/zendesktalk.svg'}
                    identifier={'zendesktalk'}
                    name={'Zendesk Talk'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'subdomain',
                            label: 'Subdomain',
                        },
                        {
                            name: 'accessKey',
                            label: 'Access Key',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'zendesksell',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/zendesksell.svg'}
                    identifier={'zendesksell'}
                    name={'Zendesk Sell'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'zendesksunshine',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/zendesksunshine.svg'}
                    identifier={'zendesksunshine'}
                    name={'Zendesk Sunshine'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'subdomain',
                            label: 'Subdomain',
                        },
                        {
                            name: 'apiToken',
                            label: 'API Token',
                        },
                        {
                            name: 'email',
                            label: 'Email',
                        },
                    ]}
                />
            ),
        },
        {
            key: 'zenefits',
            state: 'INACTIVE',
            template: (data: any) => (
                <SettingsIntegrationItem
                    icon={'logos/zenefits.svg'}
                    identifier={'zenefits'}
                    name={'Zenefits'}
                    state={data.state}
                    settingsChanged={() => {
                        reloadRef.current = !reloadRef.current;
                        setReload(reloadRef.current);
                    }}
                    fields={[
                        {
                            name: 'token',
                            label: 'Token',
                        },
                    ]}
                />
            ),
        },
    ]);

    const [integrationsDisplayed, setIntegrationsDisplayed] = useState([] as any);

    useEffect(() => {
        setLoading(true);
        GetSettings()
            .then((data: any) => {
                const map = integrations.map((integration) => {
                    return {
                        ...integration,
                        state: data[integration.key]?.state ?? 'INACTIVE',
                    };
                    return integration;
                });

                setIntegrations(map);
                setIntegrationsDisplayed(map);

                setLoading(false);
            })
            .catch((reason: any) => {
                toast.error(
                    'There was a problem on our side and we cannot load settings data at the moment,  we are doing our best to solve it! ',
                );
            });
    }, [reload]);

    const [searchTerm, setSearchTerm] = useState('');
    const handleFilterResults = (value: string) => {
        setSearchTerm(value);
        setIntegrationsDisplayed(
            integrations.filter((integration: any) =>
                integration.key.toLowerCase().includes(value.toLowerCase()),
            ),
        );
    };



  return (
    <>
        <Card
            flex='3'
            h='calc(100vh - 1rem)'
            bg='#FCFCFC'
            borderRadius='2xl'
            flexDirection='column'
            boxShadow='none'
            position='relative'
            background='gray.25'
            minWidth={609}
        >
            <CardHeader px={6} pb={2}>
                <Heading as='h1' fontSize='lg' color='gray.700'>
                    <b>Data Integrations</b>
                </Heading>
            </CardHeader>
            <CardBody>
                <div>
                    <DebouncedInput
                        className={'wfull'}
                        minLength={2}
                        onChange={(event) => handleFilterResults(event.target.value)}
                        placeholder={'Search ...'}
                        value={searchTerm}
                    >
                        <Search />
                    </DebouncedInput>

                    <div>
                        <h2 style={{ marginTop: '20px' }}>Active integrations</h2>
                        {loading && (
                            <>
                                <div style={{ marginTop: '20px' }}>
                                    <div>
                                        <Skeleton height='30px' width='100%' />
                                    </div>
                                    <div>
                                        <Skeleton height='20px' width='90%' />
                                    </div>
                                </div>
                                <div style={{ marginTop: '20px' }}>
                                    <div>
                                        <Skeleton height='30px' width='100%' />
                                    </div>
                                    <div>
                                        <Skeleton height='20px' width='90%' />
                                    </div>
                                </div>
                            </>
                        )}
                        {!loading && (
                            <>
                                {integrationsDisplayed
                                    .filter((integration: any) => integration.state === 'ACTIVE')
                                    .map((integration: any) => {
                                        return integration.template(integration);
                                    })}
                            </>
                        )}

                        <h2 style={{ marginTop: '20px' }}>Inactive integrations</h2>
                        {loading && (
                            <>
                                <div style={{ marginTop: '20px' }}>
                                    <div>
                                        <Skeleton height='30px' width='100%' />
                                    </div>
                                    <div>
                                        <Skeleton height='20px' width='90%' />
                                    </div>
                                </div>
                                <div style={{ marginTop: '20px' }}>
                                    <div>
                                        <Skeleton height='30px' width='100%' />
                                    </div>
                                    <div>
                                        <Skeleton height='20px' width='90%' />
                                    </div>
                                </div>
                            </>
                        )}
                        {!loading && (
                            <>
                                {integrationsDisplayed
                                    .filter(
                                        (integration: any) => integration.state === 'INACTIVE',
                                    )
                                    .map((integration: any) => {
                                        return integration.template(integration);
                                    })}
                            </>
                        )}
                    </div>
                </div>
            </CardBody>
        </Card>
    </>
  );
};
