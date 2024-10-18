import type { RootStore } from '@store/root';

import set from 'lodash/set';
import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { MarkerType } from '@xyflow/react';
import { Transport } from '@store/transport';
import { FlowService } from '@store/Flows/__service__';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { makeAutoSyncableGroup } from '@store/group-store';
import { FlowContactStore } from '@store/FlowContacts/FlowContact.store.ts';

import { uuidv4 } from '@utils/generateUuid.ts';
import {
  Flow,
  Contact,
  DataSource,
  FlowStatus,
  FlowContact,
  FlowParticipantStatus,
} from '@graphql/types';

export class FlowStore implements Store<Flow> {
  value: Flow = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Flow>();
  update = makeAutoSyncable.update<Flow>();
  private service: FlowService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Flow',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
    this.service = FlowService.getInstance(transport);
  }

  get id() {
    return this.value.metadata?.id;
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;

    match(path)
      .with(['status', ...P.array()], () => {
        this.service.changeStatus({
          id: this.id,
          status: this.value.status as FlowStatus,
        });
      })
      .with(['name', ...P.array()], () => {
        this.updateFlow({ nodes: this.value.nodes, edges: this.value.edges });
      });
  }

  get parsedNodes() {
    try {
      return JSON.parse(this.value.nodes);
    } catch (error) {
      console.error('Error parsing nodes:', error);

      return initialNodes;
    }
  }

  get parsedEdges() {
    try {
      return JSON.parse(this.value.edges);
    } catch (error) {
      console.error('Error parsing edges:', error);

      return initialEdges;
    }
  }

  public async updateFlow(
    { nodes, edges }: { nodes: string; edges: string },
    options?: { onError: () => void; onSuccess: () => void },
  ) {
    this.isLoading = true;

    try {
      const { flow_Merge } = await this.service.mergeFlow({
        input: {
          id: this.id,
          name: this.value.name,
          nodes,
          edges,
        },
      });

      runInAction(() => {
        this.value.nodes = flow_Merge?.nodes ?? '[]';
        this.value.edges = flow_Merge?.edges ?? '[]';
        this.root.ui.toastSuccess(
          `${this.value.name} saved`,
          `update-flow-success-${this.id}`,
        );
        setTimeout(() => {
          if (options?.onSuccess) {
            options.onSuccess();
          }
        }, 0);
      });
    } catch (e) {
      runInAction(() => {
        if (options?.onError) {
          options.onError();
        }
        this.root.ui.toastError(
          "We couldn't update the flow",
          'update-flow-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  invalidate() {
    // todo
    return Promise.resolve();
  }

  init(data: Flow) {
    const output = merge(this.value, data);

    const flowContacts = data.contacts?.map((item) => {
      this.root.flowContacts.load([item]);

      if (!item.metadata.id) {
        return;
      }

      return this.root.flowContacts.value.get(item.metadata.id)?.value;
    });
    const flowSenders = data.senders?.map((item) => {
      this.root.flowSenders.load([item]);

      return this.root.flowSenders.value.get(item?.metadata?.id)?.value;
    });

    flowContacts && set(output, 'contacts', flowContacts);
    flowSenders && set(output, 'senders', flowSenders);

    return output;
  }

  public linkContact = async (contactId: string) => {
    this.isLoading = true;

    try {
      const contactStore = this.root.contacts.value.get(contactId);

      if (contactStore?.flow) {
        await contactStore.deleteFlowContact();
      }

      const { flowContact_Add } = await this.service.addContact({
        contactId,
        flowId: this.id,
      });

      runInAction(() => {
        contactStore?.update(
          (c) => {
            c.flows = [{ ...this.value }];

            return c;
          },
          { mutate: false },
        );

        const newFlowContactValue = {
          ...flowContact_Add,
          contact: {
            id: contactId,
            metadata: {
              id: contactId,
              source: DataSource.Openline,
              appSource: DataSource.Openline,
              created: new Date().toISOString(),
              lastUpdated: new Date().toISOString(),
              sourceOfTruth: DataSource.Openline,
            },
          } as Contact,
        } as FlowContact;

        this.value.contacts = [...this.value.contacts, newFlowContactValue];

        const newFLowContact = new FlowContactStore(this.root, this.transport);

        newFLowContact.value = newFlowContactValue;
        this.root.flowContacts.value.set(
          newFlowContactValue.metadata.id,
          newFLowContact,
        );

        this.root.ui.toastSuccess(
          `Contact added to '${this.value.name}'`,
          'link-contact-to-flows-success',
        );
        contactStore?.invalidate();

        // TODO Refactor https://linear.app/customer-os/issue/COS-4820/handle-single-flow-load-invalidate-singular-flow
        // Invalidate singular flow when this issue is merged https://linear.app/customer-os/issue/COS-4823/create-getflowbyididid-query
        this.root.flows.invalidate();
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          "We couldn't add a contact to a flow",
          'link-contact-to-flows-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  public linkContacts = async (contactIds: string[]) => {
    this.isLoading = true;

    try {
      const contactStores = contactIds.map((e) => {
        return this.root.contacts.value.get(e);
      });

      await Promise.all(
        contactStores
          .filter((e): e is NonNullable<typeof e> => !!e && !!e.flowContact)
          .map((e) => e?.flowContact?.deleteFlowContact()),
      );
      await this.service.addContactBulk({
        contactId: contactIds,
        flowId: this.id,
      });

      runInAction(() => {
        contactStores.map((e) => {
          e?.update(
            (c) => {
              c.flows = [{ ...this.value }];

              return c;
            },
            { mutate: false },
          );

          return e;
        });

        this.value.contacts = [
          ...this.value.contacts,
          ...(contactStores || []).map((cs) => ({
            metadata: {
              id: uuidv4(),
              source: DataSource.Openline,
              appSource: DataSource.Openline,
              created: new Date().toISOString(),
              lastUpdated: new Date().toISOString(),
              sourceOfTruth: DataSource.Openline,
            },
            status: FlowParticipantStatus.Scheduled,
            scheduledAction: '',
            scheduledAt: new Date().toISOString(),
            contact: {
              id: cs?.id,
              metadata: {
                id: cs?.id,
                source: DataSource.Openline,
                appSource: DataSource.Openline,
                created: new Date().toISOString(),
                lastUpdated: new Date().toISOString(),
                sourceOfTruth: DataSource.Openline,
              },
            },
          })),
        ] as FlowContact[];

        this.root.ui.toastSuccess(
          `${contactIds.length} contacts added to '${this.value.name}'`,
          'link-contacts-to-flows-success',
        );
        this.root.contacts.sync({ action: 'INVALIDATE', ids: contactIds });
        this.root.flows.sync({ action: 'INVALIDATE', ids: [this.id] });
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          "We couldn't add contacts to a flow",
          'link-contacts-to-flows-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
}

const getDefaultValue = (): Flow => ({
  name: '',
  status: FlowStatus.Inactive,
  description: '',
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
  statistics: {
    total: 0,
    pending: 0,
    completed: 0,
    goalAchieved: 0,
  },
  contacts: [],
  senders: [],
  nodes: JSON.stringify(initialNodes),
  edges: JSON.stringify(initialEdges),
});
const initialNodes = [
  {
    $H: 497,
    data: { action: 'FLOW_START', entity: null, triggerType: null },
    height: 48,
    id: 'tn-1',
    internalId: '00370d94-5f6d-4d00-a1c0-3147413da9fb',
    measured: { height: 48, width: 300 },
    position: { x: 12, y: 12 },
    properties: { 'org.eclipse.elk.portConstraints': 'FIXED_ORDER' },
    sourcePosition: 'bottom',
    targetPosition: 'top',
    type: 'trigger',
    width: 300,
    x: 12,
    y: 12,
  },
  {
    $H: 499,
    data: { action: 'FLOW_END' },
    height: 48,
    id: 'tn-2',
    internalId: 'ba2070b8-79ad-4f59-8b5a-c4dd77c8cff5',
    measured: { height: 48, width: 131 },
    position: { x: 96.5, y: 160 },
    properties: { 'org.eclipse.elk.portConstraints': 'FIXED_ORDER' },
    sourcePosition: 'bottom',
    targetPosition: 'top',
    type: 'control',
    width: 131,
    x: 96.5,
    y: 160,
  },
];

const initialEdges = [
  {
    id: 'e1-2',
    source: 'tn-1',
    target: 'tn-2',
    selected: false,
    selectable: true,
    focusable: true,
    interactionWidth: 60,
    markerEnd: {
      type: MarkerType.Arrow,
      width: 20,
      height: 20,
    },
    type: 'baseEdge',
    data: { isHovered: false },
    sections: [
      {
        id: 'e1-2_s0',
        startPoint: { x: 162, y: 60 },
        endPoint: { x: 162, y: 160 },
        incomingShape: 'tn-1',
        outgoingShape: 'tn-2',
      },
    ],
    container: 'root',
  },
];
