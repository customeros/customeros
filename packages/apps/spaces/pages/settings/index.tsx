import type { NextPage } from 'next';
import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import {
  DeleteHubspotSettings,
  DeleteJiraSettings,
  DeleteSmartsheetSettings,
  DeleteTrelloSettings,
  DeleteZendeskSettings,
  GetSettings,
  UpdateHubspotSettings,
  UpdateJiraSettings,
  UpdateSmartsheetSettings,
  UpdateTrelloSettings,
  UpdateZendeskSettings,
  HubspotSettings,
  Settings,
} from '../../services';
import { Button } from '../../components';
import styles from './settings.module.scss';
import { ArrowLeft } from '../../components/ui-kit/atoms';
import { useRouter } from 'next/router';

const Settings: NextPage = () => {
  const router = useRouter();
  const [settings, setSettingsExist] = useState<Settings>({
    zendeskExists: false,
    smartSheetExists: false,
    hubspotExists: false,
    jiraExists: false,
    trelloExists: false,
  });

  const [hubspotPrivateAppKey, setHubspotPrivateAppKey] = useState<string>('');

  const [zendeskAPIKey, setZendeskApiKey] = useState<string>('');
  const [zendeskSubdomain, setZendeskSubdomain] = useState<string>('');
  const [zendeskAdminEmail, setZendeskAdminEmail] = useState<string>('');

  const [smartSheetId, setSmartsheetId] = useState<string>('');
  const [smartSheetAccessToken, setSmartsheetAccessToken] =
    useState<string>('');

  const [jiraAPIToken, setJiraAPIToken] = useState<string>('');
  const [jiraDomain, setJiraDomain] = useState<string>('');
  const [jiraEmail, setJiraEmail] = useState<string>('');

  const [trelloToken, setTrelloToken] = useState<string>('');

  useEffect(() => {
    GetSettings()
      .then((data: Settings) => {
        setSettingsExist(data);
      })
      .catch((reason: any) => {
        toast.error(
          'There was a problem on our side and we cannot load settings data at the moment,  we are doing our best to solve it! ',
        );
      });
  }, []);

  const resetZendesk = () => {
    setZendeskApiKey('');
    setZendeskSubdomain('');
    setZendeskAdminEmail('');
  };

  const resetHubspot = () => {
    setHubspotPrivateAppKey('');
  };

  const resetSmartsheet = () => {
    setSmartsheetId('');
    setSmartsheetAccessToken('');
  };

  const resetJira = () => {
    setJiraAPIToken('');
    setJiraDomain('');
    setJiraEmail('');
  };

  const resetTrello = () => {
    setTrelloToken('');
  };

  const handleSubmitHubspotSettings = () => {
    UpdateHubspotSettings({ hubspotPrivateAppKey })
      .then(() => {
        toast.success('Settings updated successfully!');
        setSettingsExist({
          ...settings,
          hubspotExists: true,
        });
        resetHubspot();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };

  const handleSubmitZendeskSettings = () => {
    UpdateZendeskSettings({
      zendeskSubdomain,
      zendeskAdminEmail,
      zendeskAPIKey,
    })
      .then(() => {
        toast.success('Settings updated successfully!');
        setSettingsExist({
          ...settings,
          zendeskExists: true,
        });
        resetZendesk();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleSubmitSmartsheetSettings = () => {
    UpdateSmartsheetSettings({
      smartSheetAccessToken,
      smartSheetId,
    })
      .then(() => {
        toast.success('Settings updated successfully!');
        setSettingsExist({
          ...settings,
          smartSheetExists: true,
        });
        resetSmartsheet();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleSubmitJiraSettings = () => {
    UpdateJiraSettings({
      jiraAPIToken,
      jiraDomain,
      jiraEmail,
    })
      .then(() => {
        toast.success('Settings updated successfully!');
        setSettingsExist({
          ...settings,
          jiraExists: true,
        });
        resetJira();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleSubmitTrelloSettings = () => {
    UpdateTrelloSettings({
      trelloToken,
    })
      .then(() => {
        toast.success('Settings updated successfully!');
        setSettingsExist({
          ...settings,
          trelloExists: true,
        });
        resetJira();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };

  const handleDeleteHubspot = () => {
    DeleteHubspotSettings()
      .then(() => {
        setSettingsExist({
          ...settings,
          hubspotExists: false,
        });
        resetHubspot();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleDeleteZendesk = () => {
    DeleteZendeskSettings()
      .then(() => {
        setSettingsExist({
          ...settings,
          zendeskExists: false,
        });
        resetZendesk();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleDeleteSmartsheetSettings = () => {
    DeleteSmartsheetSettings()
      .then(() => {
        setSettingsExist({
          ...settings,
          smartSheetExists: false,
        });
        resetSmartsheet();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleDeleteJiraSettings = () => {
    DeleteJiraSettings()
      .then(() => {
        setSettingsExist({
          ...settings,
          jiraExists: false,
        });
        resetJira();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };
  const handleDeleteTrelloSettings = () => {
    DeleteTrelloSettings()
      .then(() => {
        setSettingsExist({
          ...settings,
          trelloExists: false,
        });
        resetTrello();
      })
      .catch(() => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };

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
        <h1 className={styles.mainHeading}>Settings</h1>
      </div>

      <div className={styles.settingsContainer}>
        <article className={styles.gridItem}>
          <h2 className={styles.heading}>Hubspot</h2>

          <label htmlFor='openline-hubspot-api-key' className={styles.label}>
            API key
          </label>
          <input
            value={
              settings.hubspotExists ? '************' : hubspotPrivateAppKey
            }
            disabled={settings.hubspotExists}
            className={styles.input}
            onChange={({ target: { value } }) => setHubspotPrivateAppKey(value)}
          />
          <div className={styles.buttonSection}>
            {settings.hubspotExists ? (
              <Button onClick={handleDeleteHubspot} mode='danger'>
                Revoke
              </Button>
            ) : (
              <Button onClick={handleSubmitHubspotSettings} mode='primary'>
                Save
              </Button>
            )}
          </div>
        </article>

        <article className={styles.gridItem}>
          <h2 className={styles.heading}>Zendesk</h2>
          <label htmlFor='openline-zendesk-api-key' className={styles.label}>
            API key
          </label>
          <input
            value={settings.zendeskExists ? '*************' : zendeskAPIKey}
            id='openline-zendesk-api-key'
            disabled={settings.zendeskExists}
            className={styles.input}
            onChange={({ target: { value } }) => setZendeskApiKey(value)}
          />
          <label htmlFor='openline-zendesk-subdomain' className={styles.label}>
            Subdomain
          </label>
          <input
            value={settings.zendeskExists ? '*************' : zendeskSubdomain}
            id='openline-zendesk-subdomain'
            disabled={settings.zendeskExists}
            className={styles.input}
            onChange={({ target: { value } }) => setZendeskSubdomain(value)}
          />
          <label
            htmlFor='openline-zendesk-admin-email'
            className={styles.label}
          >
            Admin Email
          </label>
          <input
            value={settings.zendeskExists ? '*************' : zendeskAdminEmail}
            id='openline-zendesk-admin-email'
            disabled={settings.zendeskExists}
            className={styles.input}
            onChange={({ target: { value } }) => setZendeskAdminEmail(value)}
          />
          <div className={styles.buttonSection}>
            {settings.zendeskExists ? (
              <Button onClick={handleDeleteZendesk} mode='danger'>
                Revoke
              </Button>
            ) : (
              <Button onClick={handleSubmitZendeskSettings} mode='primary'>
                Save
              </Button>
            )}
          </div>
        </article>

        <article className={styles.gridItem}>
          <h2 className={styles.heading}>Smartsheet</h2>
          <label htmlFor='openline-smartsheet-id' className={styles.label}>
            ID
          </label>
          <input
            value={
              settings.smartSheetExists ? '******************' : smartSheetId
            }
            id='openline-zendesk-api-key'
            disabled={settings.smartSheetExists}
            className={styles.input}
            onChange={({ target: { value } }) => setSmartsheetId(value)}
          />
          <label htmlFor='openline-smartsheet-api-key' className={styles.label}>
            API key
          </label>
          <input
            value={
              settings.smartSheetExists
                ? '******************'
                : smartSheetAccessToken
            }
            id='openline-zendesk-api-key'
            className={styles.input}
            disabled={settings.smartSheetExists}
            onChange={({ target: { value } }) =>
              setSmartsheetAccessToken(value)
            }
          />
          <div className={styles.buttonSection}>
            {settings.smartSheetExists ? (
              <Button onClick={handleDeleteSmartsheetSettings} mode='danger'>
                Revoke
              </Button>
            ) : (
              <Button onClick={handleSubmitSmartsheetSettings} mode='primary'>
                Save
              </Button>
            )}
          </div>
        </article>

        <article className={styles.gridItem}>
          <h2 className={styles.heading}>Jira</h2>
          <label className={styles.label}>API Token</label>
          <input
            value={settings.jiraExists ? '******************' : jiraAPIToken}
            disabled={settings.jiraExists}
            className={styles.input}
            onChange={({ target: { value } }) => setJiraAPIToken(value)}
          />

          <label className={styles.label}>Domain</label>
          <input
            value={settings.jiraExists ? '******************' : jiraDomain}
            className={styles.input}
            disabled={settings.jiraExists}
            onChange={({ target: { value } }) => setJiraDomain(value)}
          />

          <label className={styles.label}>Email</label>
          <input
            value={settings.jiraExists ? '******************' : jiraEmail}
            className={styles.input}
            disabled={settings.jiraExists}
            onChange={({ target: { value } }) => setJiraEmail(value)}
          />

          <div className={styles.buttonSection}>
            {settings.jiraExists ? (
              <Button onClick={handleDeleteJiraSettings} mode='danger'>
                Revoke
              </Button>
            ) : (
              <Button onClick={handleSubmitJiraSettings} mode='primary'>
                Save
              </Button>
            )}
          </div>
        </article>

        <article className={styles.gridItem}>
          <h2 className={styles.heading}>Trello</h2>

          <label className={styles.label}>Token</label>
          <input
            value={settings.trelloExists ? '******************' : trelloToken}
            className={styles.input}
            disabled={settings.trelloExists}
            onChange={({ target: { value } }) => setTrelloToken(value)}
          />

          <div className={styles.buttonSection}>
            {settings.trelloExists ? (
              <Button onClick={handleDeleteTrelloSettings} mode='danger'>
                Revoke
              </Button>
            ) : (
              <Button onClick={handleSubmitTrelloSettings} mode='primary'>
                Save
              </Button>
            )}
          </div>
        </article>
      </div>
    </div>
  );
};

export default Settings;
