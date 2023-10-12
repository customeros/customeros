import axios from 'axios';

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
      .then(({ data, error }: any) => {
        if (data) {
          resolve(data);
        } else {
          reject(error);
        }
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
      .then(({ data, error }: any) => {
        if (data) {
          resolve(data);
        } else {
          reject(error);
        }
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}

export function GetIntegrationsSettings(): Promise<any> {
  return new Promise((resolve, reject) =>
    axios
      .get('/sa/integrations')
      .then(({ data, error }: any) => {
        if (data) {
          resolve(data);
        } else {
          reject(error);
        }
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}

export function UpdateIntegrationSettings(
  identifier: string,
  data: any,
): Promise<any> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/integration`, {
        [identifier]: data,
      })
      .then((response: any) => {
        if (response.data) {
          resolve(response.data);
        } else {
          reject(response.error);
        }
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}

export function DeleteIntegrationSettings(identifier: string): Promise<any> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/integration/${identifier}`)
      .then((response: any) => {
        if (response.data) {
          resolve(response.data);
        } else {
          reject(response.error);
        }
      })
      .catch((reason) => {
        reject(reason);
      }),
  );
}
