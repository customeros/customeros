import { fromZonedTime } from 'date-fns-tz';
import { Persister } from '@store/persister';
import { runInAction, makeAutoObservable } from 'mobx';

type NetworkStatus = 'offline' | 'online';

export class WindowManager {
  lastActiveAt: Date | null = null;
  networkStatus: NetworkStatus = navigator?.onLine ? 'online' : 'offline';
  private persister = Persister.getSharedInstance('Session');

  constructor() {
    makeAutoObservable(this);

    this.hydrateLastActiveAt();

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

  private async hydrateLastActiveAt() {
    try {
      const loadedTimestamp = await this.persister?.getItem<number>(
        'lastActiveAt',
      );

      if (loadedTimestamp) {
        this.lastActiveAt = new Date(loadedTimestamp);
      }
    } catch (e) {
      console.error('Failed to hydrate lastActiveAt', e);
    }
  }

  private async persistLastActiveAt() {
    if (this.networkStatus === 'offline') return;

    runInAction(() => {
      this.lastActiveAt = new Date();
    });

    try {
      if (this.lastActiveAt) {
        await this.persister?.setItem(
          'lastActiveAt',
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
