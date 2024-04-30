import axios from 'axios';

// TODO: remove useless wrapping promises from those functions. Axios already returns a Promise.
export interface OAuthUserSettingsInterface {
  gmailSyncEnabled: boolean;
  googleCalendarSyncEnabled: boolean;
}

export interface SlackSettingsInterface {
  slackEnabled: boolean;
}

export function GetGoogleSettings(
  identifier: string,
): Promise<OAuthUserSettingsInterface> {
  return new Promise((resolve, reject) =>
    axios
      .get(`/sa/user/settings/google/${identifier}`)
      .then(({ data }) => {
        resolve(data);
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}

export function GetSlackSettings(): Promise<SlackSettingsInterface> {
  return new Promise((resolve, reject) =>
    axios
      .get(`/sa/user/settings/slack`)
      .then(({ data }) => {
        resolve(data);
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}

export function GetIntegrationsSettings(): Promise<unknown> {
  return new Promise((resolve, reject) =>
    axios
      .get('/sa/integrations')
      .then(({ data }) => {
        resolve(data);
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}

export function UpdateIntegrationSettings(
  identifier: string,
  data: unknown,
): Promise<unknown> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/integration`, {
        [identifier]: data,
      })
      .then(({ data }) => {
        resolve(data);
      })
      .catch((error) => {
        reject(error);
      }),
  );
}

export function DeleteIntegrationSettings(
  identifier: string,
): Promise<unknown> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/integration/${identifier}`)
      .then(({ data }) => {
        resolve(data);
      })
      .catch((error) => {
        reject(error);
      }),
  );
}
