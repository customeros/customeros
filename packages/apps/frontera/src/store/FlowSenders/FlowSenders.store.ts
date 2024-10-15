import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { FlowSendersService } from '@store/FlowSenders/__service__';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowSenderStore } from '@store/FlowSenders/FlowSender.store.ts';

import { FlowSender } from '@graphql/types';

export class FlowSendersStore implements GroupStore<FlowSender> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<FlowSender>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<FlowSender>();
  totalElements = 0;
  private service: FlowSendersService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'FlowSenders',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: FlowSenderStore,
    });
    this.service = FlowSendersService.getInstance(transport);
  }

  createFlowSender = async (flowId: string, userId: string) => {
    this.isLoading = true;

    const flow = this.root.flows.value.get(flowId);
    const user = this.root.users.value.get(userId);

    const newFlowSender = new FlowSenderStore(this.root, this.transport);
    const tempId = newFlowSender.id;

    newFlowSender.value.flow = flow?.value;
    newFlowSender.value.user = user?.value;

    let serverId = '';

    this.value.set(tempId, newFlowSender);
    this.isLoading = true;

    try {
      const { flowSender_Merge } = await this.service.createFlowSender({
        flowId,
        input: {
          userId,
        },
      });

      runInAction(() => {
        serverId = flowSender_Merge.metadata.id;
        newFlowSender.setId(serverId);

        this.value.set(serverId, newFlowSender);
        this.value.delete(tempId);
        flow?.value.senders.push(newFlowSender.value);

        this.sync({
          action: 'APPEND',
          ids: [serverId],
        });
        this.root.flows.sync({
          action: 'INVALIDATE',
          ids: [flowId],
        });
        this.root.ui.toastSuccess(
          'Sender added to flow',
          'link-sender-to-flow',
        );
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
        this.root.ui.toastError(
          "We couldn't add a sender to a flow",
          'link-sender-to-flow-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  deleteFlowSender = async (senderId: string, flowId: string) => {
    this.isLoading = true;

    const flow = this.root.flows.value.get(flowId);
    let prevSenders: Array<FlowSender> = [];

    if (flow) {
      prevSenders = flow.value.senders;
      flow.value.senders = flow.value.senders.filter(
        (e) => e.metadata.id !== senderId,
      );
    }

    this.isLoading = true;

    try {
      const { flowSender_Delete } = await this.service.deleteFlowSender({
        id: senderId,
      });

      runInAction(() => {
        this.sync({
          action: 'DELETE',
          ids: [senderId],
        });
        this.root.flows.sync({
          action: 'INVALIDATE',
          ids: [flowId],
        });

        if (!flowSender_Delete.result) {
          if (flow) {
            flow.value.senders = prevSenders;
          }
        }
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;

        this.root.ui.toastError(
          "We couldn't remove sender from a flow",
          'link-sender-to-flow-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
}
