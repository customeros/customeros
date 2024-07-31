import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { LogEntry } from '@graphql/types';

import { LogEntryStore } from './LogEntry.store';
import { LogEntriesService } from './__service__/LogEntries.service';

export class LogEntriesStore implements GroupStore<LogEntry> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, LogEntryStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<LogEntry>();
  subscribe = makeAutoSyncableGroup.subscribe;
  private service: LogEntriesService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = LogEntriesService.getInstance(transport);

    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'LogEntries',
      getItemId: (item) => item.id,
      ItemStore: LogEntryStore,
    });
  }

  async create(
    organizationId: string,
    { content = '', tags }: { tags: string[]; content: string },
  ) {
    const logEntry = new LogEntryStore(this.root, this.transport);

    logEntry.value.content = content;
    this.value.set(logEntry.id, logEntry);
    this.root.timelineEvents.value.get(organizationId)?.push(logEntry);

    let serverId = '';

    try {
      this.isLoading = true;

      const { logEntry_CreateForOrganization } =
        await this.service.createLogEntry({
          organizationId,
          logEntry: { content, tags: tags.map((tag) => ({ name: tag })) },
        });

      runInAction(() => {
        serverId = logEntry_CreateForOrganization;

        this.value.delete(logEntry.id);
        logEntry.id = serverId;

        this.value.set(serverId, logEntry);

        const timeline = this.root.timelineEvents.value.get(organizationId);

        timeline?.splice(timeline.length - 1, 1, logEntry);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        setTimeout(() => {
          logEntry.invalidate();
          this.sync({ action: 'APPEND', ids: [serverId] });
          this.root.timelineEvents.invalidateTimeline(organizationId);
        }, 500);
      });
    }
  }
}
