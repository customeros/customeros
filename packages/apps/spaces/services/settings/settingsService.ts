import { HubspotSettings, SmartsheetSettings, ZendeskSettings } from './types';
import axios from 'axios';
export function GetSettings(): Promise<any> {
  return new Promise((resolve, reject) =>
    axios
      .get('/sa/settings/')
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
export function UpdateHubspotSettings(
  data: HubspotSettings,
): Promise<HubspotSettings> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/settings/hubspot`, data)
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
export function UpdateZendeskSettings(
  data: ZendeskSettings,
): Promise<ZendeskSettings> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/settings/zendesk`, data)
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

export function UpdateSmartsheetSettings(
  data: SmartsheetSettings,
): Promise<SmartsheetSettings> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/settings/smartSheet`, data)
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

export function UpdateJiraSettings(data: any): Promise<any> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/settings/jira`, data)
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

export function UpdateTrelloSettings(data: any): Promise<any> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/sa/settings/trello`, data)
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

export function DeleteSmartsheetSettings(): Promise<SmartsheetSettings> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/settings/smartSheet`)
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
export function DeleteHubspotSettings(): Promise<SmartsheetSettings> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/settings/hubspot`)
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

export function DeleteZendeskSettings(): Promise<SmartsheetSettings> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/settings/zendesk`)
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

export function DeleteJiraSettings(): Promise<SmartsheetSettings> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/settings/jira`)
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

export function DeleteTrelloSettings(): Promise<SmartsheetSettings> {
  return new Promise((resolve, reject) =>
    axios
      .delete(`/sa/settings/trello`)
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
