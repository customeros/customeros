import type { NextPage } from 'next';
import React, { useEffect, useRef, useState } from 'react';
import { toast } from 'react-toastify';
import { GetSettings } from '../../services';
import { Button } from '../../components';
import styles from './settings.module.scss';
import { ArrowLeft } from '../../components/ui-kit/atoms';
import { useRouter } from 'next/router';
import { Skeleton } from '../../components/ui-kit/atoms/skeleton';
import { SettingsIntegrationItem } from '../../components/ui-kit/molecules/settings-integration-item';

const Settings: NextPage = () => {
  const router = useRouter();

  const [reload, setReload] = useState<boolean>(false);
  const reloadRef = useRef<boolean>(reload);

  const [loading, setLoading] = useState<boolean>(true);

  //states: ACTIVE OR INACTIVE
  //TODO: switch to a different state when the integration is being configured to fetch running and error states
  const [integrations, setIntegrations] = useState([
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
      key: 'zendesk',
      state: 'INACTIVE',
      template: (data: any) => (
        <SettingsIntegrationItem
          icon={'/logos/zendesk.svg'}
          identifier={'zendesk'}
          name={'Zendesk'}
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
      key: 'baton',
      state: 'INACTIVE',
      template: (data: any) => (
        <SettingsIntegrationItem
          icon={'logos/openline_small.svg'}
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
          icon={'logos/openline_small.svg'}
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
  ]);

  useEffect(() => {
    setLoading(true);
    GetSettings()
      .then((data: any) => {
        setIntegrations(
          integrations.map((integration) => {
            return {
              ...integration,
              state: data[integration.key]?.state ?? 'INACTIVE'
            };
            return integration;
          }),
        );

        setLoading(false);
      })
      .catch((reason: any) => {
        toast.error(
          'There was a problem on our side and we cannot load settings data at the moment,  we are doing our best to solve it! ',
        );
      });
  }, [reload]);

  return (
    <div className={styles.pageContainer}>
      <div className={styles.headingSection}>
        <Button
          mode='secondary'
          icon={<ArrowLeft />}
          onClick={() => router.back()}
        >
          Back
        </Button>
      </div>

      <div className={styles.settingsContainer}>
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
            <h2 style={{ marginTop: '20px' }}>Active integrations</h2>
            {integrations
              .filter((integration) => integration.state === 'ACTIVE')
              .map((integration) => {
                return integration.template(integration);
              })}
          </>
        )}

        {!loading && (
          <>
            <h2 style={{ marginTop: '20px' }}>Inactive integrations</h2>
            {integrations
              .filter((integration) => integration.state === 'INACTIVE')
              .map((integration) => {
                return integration.template(integration);
              })}
          </>
        )}
      </div>
    </div>
  );
};

export default Settings;
