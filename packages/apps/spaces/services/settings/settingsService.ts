import axios from 'axios';
import {env} from "string-env-interpolation";

export interface OAuthUserSettingsInterface {
    gmailSyncEnabled: boolean;
    googleCalendarSyncEnabled: boolean;
}

export function GetOAuthUserSettings(identifier:string): Promise<OAuthUserSettingsInterface> {
    return new Promise((resolve, reject) =>
        axios
            .get(`/sa/user/settings/oauth/${identifier}`)
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
