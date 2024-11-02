import { fromZonedTime } from 'date-fns-tz';
import { Persister } from '@store/persister';
import { when, runInAction, makeAutoObservable } from 'mobx';

import type { RootStore } from './root';

type NetworkStatus = 'offline' | 'online';

export class WindowManager {
  lastActiveAt: Date | null = null;
  networkStatus: NetworkStatus = navigator?.onLine ? 'online' : 'offline';
  private persisterKey: string = '';
  private persister = Persister.getSharedInstance('Meta');

  constructor(private root: RootStore) {
    makeAutoObservable(this);

    when(
      () => !!this.root.session?.value?.tenant,
      () => {
        this.persisterKey = `lastActiveAt-${this.root.session?.value?.tenant}`;
        this.hydrateLastActiveAt(this.persisterKey);
      },
    );

    window.addEventListener('blur', () => {
      this.persistLastActiveAt();
    });
    window?.addEventListener('online', () => {
      this.setNetworkStatus('online');
    });
    window?.addEventListener('offline', () => {
      this.persistLastActiveAt();
      this.setNetworkStatus('offline');
    });
  }

  /**
   * @returns Datetime in UTC
   * @default if no previous lastActiveAt exists it returns new Date() as UTC
   * */
  public getLastActiveAtUTC() {
    return this.lastActiveAt
      ? fromZonedTime(
          this.lastActiveAt,
          Intl.DateTimeFormat().resolvedOptions().timeZone,
        )
      : fromZonedTime(
          new Date(),
          Intl.DateTimeFormat().resolvedOptions().timeZone,
        );
  }

  public clearPersisterKey() {
    this.persisterKey = '';
  }

  private async hydrateLastActiveAt(idbKey: string) {
    try {
      const loadedTimestamp = await this.persister?.getItem<number>(idbKey);

      if (loadedTimestamp) {
        this.lastActiveAt = new Date(loadedTimestamp);
      }
    } catch (e) {
      console.error('Failed to hydrate lastActiveAt', e);
    }
  }

  private async persistLastActiveAt() {
    if (!this.persisterKey) return;
    if (this.networkStatus === 'offline') return;

    runInAction(() => {
      this.lastActiveAt = new Date();
    });

    try {
      if (this.lastActiveAt) {
        await this.persister?.setItem(
          this.persisterKey,
          this.lastActiveAt.valueOf(),
        );
      }
    } catch (e) {
      console.error('Failed persisting lastActiveAt', e);
    }
  }

  private setNetworkStatus(status: NetworkStatus) {
    this.networkStatus = status;
  }
}
