import { Channel } from 'phoenix';
import { set } from 'date-fns/set';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { addDays } from 'date-fns/addDays';
import { Transport } from '@store/transport';
import { UserStore } from '@store/Users/User.store';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import { Reminder, DataSource } from '@graphql/types';

import { RemindersService } from './__service__/Reminders.service';

export class ReminderStore implements Store<Reminder> {
  value: Reminder = getDefaultValue();
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Reminder>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Reminder>();
  private service: RemindersService;

  constructor(public root: RootStore, public transport: Transport) {
    this.value = getDefaultValue();
    this.service = RemindersService.getInstance(transport);

    makeAutoSyncable(this, {
      channelName: 'Reminder',
      mutator: this.save,
      getId: (item) => item?.metadata.id,
    });
    makeAutoObservable(this);
  }

  async invalidate() {}

  async updateReminder() {
    const { metadata, content, dismissed, dueDate } = this.value;

    try {
      this.isLoading = true;
      await this.service.updateReminder({
        input: { id: metadata.id, content, dismissed, dueDate },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async save() {
    return this.updateReminder();
  }

  get id() {
    return this.value.metadata.id;
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }
}

const getDefaultValue = (): Reminder => ({
  metadata: {
    id: crypto.randomUUID(),
    appSource: 'web',
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  content: '',
  dismissed: false,
  dueDate: set(addDays(new Date(), 1), {
    hours: 9,
    minutes: 0,
    seconds: 0,
    milliseconds: 0,
  }).toISOString(),
  owner: UserStore.getDefaultValue(),
});
