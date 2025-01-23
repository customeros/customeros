import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowParticipantsService } from '@store/FlowParticipants/__service__';
import { FlowParticipantStore } from '@store/FlowParticipants/FlowParticipant.store.ts';

import { FlowEntityType, FlowParticipant } from '@graphql/types';

export class FlowParticipantsStore implements GroupStore<FlowParticipant> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<FlowParticipant>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<FlowParticipant>();
  totalElements = 0;
  private service: FlowParticipantsService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'FlowParticipants',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: FlowParticipantStore,
    });
    this.service = FlowParticipantsService.getInstance(transport);
  }

  public addFlowParticipants = async (
    entityIds: string[],
    flowId: string,
    entityType?: FlowEntityType,
  ) => {
    return this.service.createFlowParticipants({
      entitiIds: entityIds,
      flowId,
      entityType: entityType || FlowEntityType.Contact,
    });
  };
  public addFlowParticipant = async (
    entityId: string,
    flowId: string,
    entityType?: FlowEntityType,
  ) => {
    return this.service.createFlowParticipant({
      entityId,
      flowId,
      entityType: entityType || FlowEntityType.Contact,
    });
  };

  public deleteFlowParticipants = async (ids: string[], flowId?: string) => {
    if (!ids.length) return;
    this.isLoading = true;

    const FlowParticipants = ids.map(
      (id) => this.value.get(id) as FlowParticipantStore,
    );

    const contactStores = FlowParticipants.map((fc) => fc?.contact);
    const flowStores = contactStores.flatMap((cs) => cs?.flows);

    try {
      await this.service.deleteFlowParticipants({
        id: ids,
      });

      runInAction(() => {
        contactStores.forEach((contact) => {
          contact?.draft();

          const foundFlow = contact?.value.flows.findIndex((f) => f === flowId);

          if (foundFlow && foundFlow !== -1) {
            contact?.value.flows.splice(foundFlow, 1);
          }
          contact?.commit({ syncOnly: true });
        });
        flowStores.forEach((c) => {
          c?.update(
            (c) => {
              c.participants = c.participants.filter(
                (e) => !ids.includes(e.metadata.id),
              );

              return c;
            },
            { mutate: false },
          );
        });

        this.root.contacts.sync({
          action: 'INVALIDATE',
          ids: FlowParticipants.map((e) => e.contactId),
        });

        const flowsIds = flowStores
          .flatMap((c) => c?.id)
          .filter((id): id is string => typeof id === 'string');

        this.root.flows.sync({
          action: 'INVALIDATE',
          ids: flowsIds,
        });
        this.root.ui.toastSuccess(
          `Contacts removed from flows`,
          'unlink-contact-from-flow-success',
        );
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          `We couldn't remove a contact from a flow`,
          'unlink-contact-from-flow-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
  public deleteFlowParticipant = async (
    flowId: string,
    participantId: string,
  ) => {
    this.isLoading = true;

    const contactId = this.value.get(participantId)?.value?.entityId;

    const flowName = this.root?.flows?.value.get(flowId)?.value.name;

    if (!contactId) {
      runInAction(() => {
        this.root.ui.toastError(
          `We couldn't remove a contact from a flow`,
          'unlink-contact-from-flow-error',
        );
      });

      return;
    }

    try {
      await this.service.deleteFlowParticipant({
        id: participantId,
      });

      runInAction(() => {
        this.root.contacts.value.get(contactId)?.draft();

        const foundFlow = this.root.contacts.value
          .get(contactId)
          ?.value.flows.findIndex((f) => f === flowId);

        if (foundFlow !== -1 && foundFlow) {
          this.root.contacts.value
            .get(contactId)
            ?.value.flows.splice(foundFlow, 1);
        }

        this.root.contacts.value.get(contactId)?.commit({ syncOnly: true });

        this.root.ui.toastSuccess(
          `Contact removed from '${flowName}'`,
          'unlink-contact-from-flow-success',
        );
        this.root.contacts.sync({
          action: 'INVALIDATE',
          ids: [contactId],
        });

        this.root.flows.sync({
          action: 'INVALIDATE',
          ids: [flowId],
        });
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          `We couldn't remove a contact from a flow`,
          'unlink-contact-from-flow-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
}
