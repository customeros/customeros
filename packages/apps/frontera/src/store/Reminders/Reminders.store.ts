import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Reminder } from '@graphql/types';

import { ReminderStore } from './Reminder.store';
import { RemindersService } from './__service__/Reminders.service';

export class RemindersStore implements GroupStore<Reminder> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, ReminderStore> = new Map();
  valueByOrganization: Map<string, ReminderStore[]> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<Reminder>();
  subscribe = makeAutoSyncableGroup.subscribe;
  private service: RemindersService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = RemindersService.getInstance(transport);

    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Reminders',
      getItemId: (item) => item.metadata.id,
      ItemStore: ReminderStore,
    });
  }

  async bootstrapByOrganization(organizationId: string) {
    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isBootstrapped = true;
      const { remindersForOrganization } =
        await this.service.getRemindersByOrganizationId({ organizationId });

      runInAction(() => {
        const reminders = remindersForOrganization
          .map((reminder) => {
            if (this.value.has(reminder.metadata.id)) {
              return this.value.get(reminder.metadata.id);
            }

            const reminderStore = new ReminderStore(this.root, this.transport);
            reminderStore.load(reminder as Reminder);
            this.value.set(reminder.metadata.id, reminderStore);

            return reminderStore;
          })
          .filter(Boolean) as ReminderStore[];

        this.valueByOrganization.set(organizationId, reminders);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      this.isLoading = false;
    }
  }

  async create(organizationId: string) {
    const userId = this.root.globalCache.value?.user.id ?? '';
    const newReminder = new ReminderStore(this.root, this.transport);
    const tempId = newReminder.id;

    const previousEntries = this.getByOrganizationId(organizationId);
    this.value.set(tempId, newReminder);
    this.valueByOrganization.set(organizationId, [
      ...previousEntries,
      newReminder,
    ]);

    try {
      this.isLoading = true;
      const { reminder_Create: serverId } = await this.service.createReminder({
        input: {
          content: newReminder.value.content ?? '',
          dueDate: newReminder.value.dueDate,
          organizationId,
          userId,
        },
      });

      runInAction(() => {
        newReminder.id = serverId ?? '';

        this.value.set(serverId ?? '', newReminder);
        this.value.delete(tempId);

        const previousEntries = this.getByOrganizationId(organizationId);
        this.valueByOrganization.set(organizationId, [
          ...previousEntries.filter((reminder) => reminder.id !== tempId),
          newReminder,
        ]);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        setTimeout(() => {
          this.sync({ action: 'APPEND', ids: [newReminder.id] });
          newReminder.invalidate();
        }, 500);
      });
    }
  }

  getByOrganizationId(organizationId: string) {
    return this.valueByOrganization.get(organizationId) ?? [];
  }
}
