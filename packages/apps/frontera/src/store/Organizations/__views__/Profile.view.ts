import { runInAction } from 'mobx';

import { OrganizationsStore } from '../Organizations.store';

export class ProfileView {
  constructor(private store: OrganizationsStore) {
    this.retrieveOrganizationForProfile();
  }

  async retrieveOrganizationForProfile() {
    try {
      const urlId = (() => {
        // get organization id from url if possible
        // necessary to bootstrap the targeted organization on a profile view
        const parts = window?.location?.pathname?.split('/');

        if (parts.length !== 3 && parts[parts.length - 1].length !== 36) {
          return null;
        }

        return parts[parts.length - 1];
      })();

      if (!urlId) return;

      await this.store.retrieve([urlId]);
    } catch (e) {
      runInAction(() => {
        this.store.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.store.isLoading = false;
      });
    }
  }
}
