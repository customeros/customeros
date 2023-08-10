import axios from 'axios';
import {env} from "string-env-interpolation";


export interface UserSettingsInterface {
    id: string;
    tenantName: string;
    username: string;
    googleOAuthAllScopesEnabled: boolean;
    googleOAuthUserAccessToken: string;
}

export function UpdateUserSettings(
    data: UserSettingsInterface,
): Promise<any> {
    return new Promise((resolve, reject) =>
        axios
            .post(`/sa/user/settings`, data,).then((response: any) => {
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

export function GetUserSettings(identifier:string): Promise<UserSettingsInterface> {
    return new Promise((resolve, reject) =>
        axios
            .get(`/sa/user/settings/${identifier}`)
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
