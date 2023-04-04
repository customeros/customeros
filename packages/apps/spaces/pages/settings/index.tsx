import type { NextPage } from 'next';
import React, { useEffect, useRef, useState } from 'react';
import { toast } from 'react-toastify';
import {
  DeleteHubspotSettings,
  DeleteJiraSettings,
  DeleteSmartsheetSettings,
  DeleteTrelloSettings,
  DeleteZendeskSettings,
  GetSettings,
  Settings,
  UpdateHubspotSettings,
  UpdateJiraSettings,
  UpdateSmartsheetSettings,
  UpdateTrelloSettings,
  UpdateZendeskSettings,
} from '../../services';
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
          icon={'/logos/hubspot.png'}
          name={'Hubspot'}
          state={data.state}
          onSave={UpdateHubspotSettings}
          onRevoke={DeleteHubspotSettings}
          settingsChanged={() => {
            reloadRef.current = !reloadRef.current;
            setReload(reloadRef.current);
          }}
          fields={[
            {
              name: 'hubspotPrivateAppKey',
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
          icon={'/logos/zendesk.png'}
          name={'Zendesk'}
          state={data.state}
          onSave={UpdateZendeskSettings}
          onRevoke={DeleteZendeskSettings}
          settingsChanged={() => {
            reloadRef.current = !reloadRef.current;
            setReload(reloadRef.current);
          }}
          fields={[
            {
              name: 'zendeskAPIKey',
              label: 'API key',
            },
            {
              name: 'zendeskSubdomain',
              label: 'Subdomain',
            },
            {
              name: 'zendeskAdminEmail',
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
          icon={'/logos/smartsheet.png'}
          name={'Smartsheet'}
          state={data.state}
          onSave={UpdateSmartsheetSettings}
          onRevoke={DeleteSmartsheetSettings}
          settingsChanged={() => {
            reloadRef.current = !reloadRef.current;
            setReload(reloadRef.current);
          }}
          fields={[
            {
              name: 'smartSheetId',
              label: 'ID',
            },
            {
              name: 'smartSheetAccessToken',
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
          icon={'/logos/jira.png'}
          name={'Jira'}
          state={data.state}
          onSave={UpdateJiraSettings}
          onRevoke={DeleteJiraSettings}
          settingsChanged={() => {
            reloadRef.current = !reloadRef.current;
            setReload(reloadRef.current);
          }}
          fields={[
            {
              name: 'jiraAPIToken',
              label: 'API Token',
            },
            {
              name: 'jiraDomain',
              label: 'Domain',
            },
            {
              name: 'jiraEmail',
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
          icon={'/logos/trello.png'}
          name={'Trello'}
          state={data.state}
          onSave={UpdateTrelloSettings}
          onRevoke={DeleteTrelloSettings}
          settingsChanged={() => {
            reloadRef.current = !reloadRef.current;
            setReload(reloadRef.current);
          }}
          fields={[
            {
              name: 'trelloAPIToken',
              label: 'API Token',
            },
            {
              name: 'trelloAPIKey',
              label: 'API key',
            },
          ]}
        />
      ),
    },
  ]);

  useEffect(() => {
    setLoading(true);
    GetSettings()
      .then((data: Settings) => {
        console.log(data);
        setIntegrations(
          integrations.map((integration) => {
            //todo switch to the generic solution one BE is done
            if (integration.key === 'hubspot') {
              return {
                ...integration,
                state: data.hubspotExists ? 'ACTIVE' : 'INACTIVE',
              };
            }
            if (integration.key === 'zendesk') {
              return {
                ...integration,
                state: data.zendeskExists ? 'ACTIVE' : 'INACTIVE',
              };
            }
            if (integration.key === 'smartsheet') {
              return {
                ...integration,
                state: data.smartSheetExists ? 'ACTIVE' : 'INACTIVE',
              };
            }
            if (integration.key === 'jira') {
              return {
                ...integration,
                state: data.jiraExists ? 'ACTIVE' : 'INACTIVE',
              };
            }
            if (integration.key === 'trello') {
              return {
                ...integration,
                state: data.trelloExists ? 'ACTIVE' : 'INACTIVE',
              };
            }
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
