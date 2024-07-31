import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

interface Field {
  name: string;
  label: string;
  textarea?: boolean;
}

export interface Integration {
  key: string;
  name: string;
  icon: string;
  fields: Field[];
  identifier: string;
  state: 'INACTIVE' | 'ACTIVE';
  isFromIntegrationApp?: boolean;
}

export class IntegrationsStore {
  value: Record<string, Integration> = {};
  isMutating = false;
  isBootstrapped = false;
  isBootstrapping = false;
  error: string | null = null;

  constructor(private root: RootStore, private transport: Transport) {
    makeAutoObservable(this);
  }

  get isLoading() {
    return this.isBootstrapping || this.isMutating;
  }

  async load() {
    if (this.root.demoMode) {
      this.value = mock as unknown as Record<string, Integration>;
      this.isBootstrapped = true;

      return;
    }

    try {
      runInAction(() => {
        this.isBootstrapping = true;
      });

      const { data } = await this.transport.http.get('/sa/integrations');

      runInAction(() => {
        this.value = data;
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isBootstrapping = false;
      });
    }
  }

  async update(identifier: string, payload: unknown) {
    Object.assign(this.value, {
      [identifier]: { state: 'ACTIVE' },
    });

    try {
      this.isMutating = true;
      this.transport.http.post('/sa/integration', {
        [identifier]: payload,
      });
      this.root.ui.toastSuccess(
        'Settings updated successfully!',
        'integration-settings-update',
      );
    } catch (err) {
      delete this.value[identifier];
      this.error = (err as Error).message;
      this.root.ui.toastError(
        `We couldn't update the Settings.`,
        'integration-settings-update-failed',
      );
    } finally {
      this.isMutating = false;
    }
  }

  async delete(identifier: string) {
    const integration = { ...this.value[identifier] };

    if (identifier in this.value) {
      delete this.value[identifier];
    }

    try {
      this.isMutating = true;
      this.transport.http.delete(`/sa/integration/${identifier}`);
      this.root.ui.toastSuccess(
        'Settings updated successfully!',
        'integration-settings-delete',
      );
    } catch (err) {
      this.value[identifier] = integration;
      this.error = (err as Error).message;
      this.root.ui.toastError(
        `We couldn't update the Settings.`,
        'integration-settings-delete-failed',
      );
    } finally {
      this.isMutating = false;
    }
  }
}

const mock = {
  intercom: {
    state: 'ACTIVE',
  },
  mixpanel: {
    state: 'ACTIVE',
  },
  notion: {
    state: 'ACTIVE',
  },
  slack: {
    state: 'ACTIVE',
  },
};
