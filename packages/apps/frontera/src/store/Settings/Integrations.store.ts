import type { RootStore } from '@store/root';

import { TransportLayer } from '@store/transport';
import { autorun, makeAutoObservable } from 'mobx';

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

  constructor(
    private rootStore: RootStore,
    private transportLater: TransportLayer,
  ) {
    makeAutoObservable(this, {}, { autoBind: true });
    autorun(() => {
      if (
        this.rootStore.sessionStore.isHydrated &&
        this.rootStore.sessionStore.isAuthenticated &&
        this.transportLater.isAuthenthicated
      ) {
        this.load();
      }
    });
  }

  get isLoading() {
    return !this.isBootstrapped;
  }

  async load() {
    try {
      this.isBootstrapping = true;
      const { data } = await this.transportLater.http.get('/sa/integrations');
      this.value = data;
    } catch (err) {
      this.error = (err as Error).message;
    } finally {
      this.isBootstrapping = false;
      if (!this.isBootstrapped) {
        this.isBootstrapped = true;
      }
    }
  }

  async update(identifer: string, payload: unknown) {
    try {
      this.isMutating = true;
      this.transportLater.http.post('/sa/integration', {
        [identifer]: payload,
      });
      this.rootStore.uiStore.toastSuccess(
        'Settings updated successfully!',
        'integration-settings-updated',
      );
    } catch (err) {
      this.error = (err as Error).message;
      this.rootStore.uiStore.toastError(
        `We couldn't update the Settings.`,
        'integration-settings-update-failed',
      );
    } finally {
      this.isMutating = false;
      this.load();
    }
  }

  async delete(identifier: string) {
    try {
      this.isMutating = true;
      this.transportLater.http.delete(`/sa/integration/${identifier}`);
      this.rootStore.uiStore.toastSuccess(
        'Settings updated successfully!',
        'integration-settings-updated',
      );
    } catch (err) {
      this.error = (err as Error).message;
      this.rootStore.uiStore.toastError(
        `We couldn't update the Settings.`,
        'integration-settings-update-failed',
      );
    } finally {
      this.isMutating = false;
      this.load();
    }
  }
}
