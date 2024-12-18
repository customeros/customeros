import { when, runInAction } from 'mobx';

import { OrganizationsStore } from '../Organizations.store';

export class ProfileView {
  constructor(private store: OrganizationsStore) {
    when(
      () => this.store.root.session.sessionToken !== null,
      this.retrieveOrganizationForProfile,
    );
  }

  retrieveOrganizationForProfile = async () => {
    try {
      runInAction(() => {
        this.store.isLoading = true;
        this.store.isBootstrapping = true;
      });

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
        this.store.isBootstrapped = true;
        this.store.isBootstrapping = false;
        this.store.isLoading = false;
      });
    }
  };
}
