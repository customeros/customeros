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
import { FlowParticipantStore } from '@store/FlowParticipants/FlowParticipant.store';

import { uuidv4 } from '@utils/generateUuid';
import {
  Flow,
  DataSource,
  FlowStatus,
  FlowParticipant,
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

  async saveStatus() {
    this.isLoading = true;

    try {
      await this.service.changeStatus({
        id: this.id,
        status: this.value.status as FlowStatus,
      });
    } catch (error) {
      this.root.ui.toastError(
        "We couldn't update the flow",
        'update-flow-error',
      );
    } finally {
      runInAction(() => {
        this.invalidate();
      });
      this.isLoading = false;
    }
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;

    match(path)
      .with(['status', ...P.array()], () => {
        this.saveStatus();
      })
      .with(['name', ...P.array()], () => {
        // todo COS-5311 - use another mutation to not update nodes and edges when updating the name
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

  async invalidate() {
    try {
      this.isLoading = true;

      const { flow } = await this.service.getFlow(this.id);

      this.init(flow);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  init(data: Flow) {
    const output = merge(this.value, data);

    const flowContacts = data.participants?.map((item) => {
      this.root.flowParticipants.load([item]);

      if (!item.metadata.id) {
        return;
      }

      return this.root.flowParticipants.value.get(item.metadata.id)?.value;
    });

    const flowSenders = data.senders?.map((item) => {
      this.root.flowSenders.load([item]);

      return this.root.flowSenders.value.get(item?.metadata?.id)?.value;
    });

    flowContacts && set(output, 'participants', flowContacts);
    flowSenders && set(output, 'senders', flowSenders);

    return output;
  }

  public linkContact = async (contactId: string) => {
    this.isLoading = true;

    try {
      const contactStore = this.root.contacts.value.get(contactId);

      const { flowParticipant_Add } =
        await this.root.flowParticipants.addFlowParticipant(contactId, this.id);

      runInAction(() => {
        contactStore?.update(
          (c) => {
            c.flows = [...(c.flows ?? []), { ...this.value }];

            return c;
          },
          { mutate: false },
        );

        const newFLowContact = new FlowParticipantStore(
          this.root,
          this.transport,
        );

        const newFlowContactValue = {
          metadata: {
            ...newFLowContact.value.metadata,
            id: flowParticipant_Add.metadata.id,
          },
          entityType: 'CONTACT',
          entityId: contactId,
          status: FlowParticipantStatus.Scheduled,
          executions: [],
        };

        this.value.participants = [
          ...this.value.participants,
          newFlowContactValue,
        ];
        this.value.statistics.onHold += 1;
        this.value.statistics.total += 1;

        newFLowContact.value = newFlowContactValue;
        this.root.flowParticipants.value.set(
          newFlowContactValue.metadata.id,
          newFLowContact,
        );

        this.root.ui.toastSuccess(
          `Contact added to flow`,
          'link-contact-to-flows-success',
        );
        contactStore?.invalidate();
        setTimeout(() => {
          this.invalidate();
        }, 400);
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

      await this.root.flowParticipants.addFlowParticipants(contactIds, this.id);

      runInAction(() => {
        contactStores.map((e) => {
          e?.update(
            (c) => {
              c.flows = [...(c.flows ?? []), { ...this.value }];

              return c;
            },
            { mutate: false },
          );

          return e;
        });

        this.value.participants = [
          ...this.value.participants,
          ...(contactStores || []).map((cs) => ({
            entityType: 'CONTACT',
            entityId: cs?.id,
            metadata: {
              id: uuidv4(),
              source: DataSource.Openline,
              appSource: DataSource.Openline,
              created: new Date().toISOString(),
              lastUpdated: new Date().toISOString(),
              sourceOfTruth: DataSource.Openline,
            },
          })),
        ] as FlowParticipant[];

        this.root.ui.toastSuccess(
          `${contactIds.length} contacts added to flow`,
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
  status: FlowStatus.Off,
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
    onHold: 0,
    ready: 0,
    scheduled: 0,
    inProgress: 0,
    completed: 0,
    goalAchieved: 0,
  },
  participants: [],
  // deprecated but needed for type compatibility
  senders: [],
  nodes: JSON.stringify(initialNodes),
  edges: JSON.stringify(initialEdges),
});

const initialNodes = [
  {
    $H: 497,
    data: {
      action: 'FLOW_START',
      entity: 'CONTACT',
      triggerType: 'RecordAddedManually',
    },
    height: 83,
    id: 'tn-1',
    position: { x: 12, y: 11 },
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
    height: 56,
    id: 'tn-2',
    measured: { height: 56, width: 131 },
    position: { x: 84, y: 195 },
    properties: { 'org.eclipse.elk.portConstraints': 'FIXED_ORDER' },
    sourcePosition: 'bottom',
    targetPosition: 'top',
    type: 'control',
    width: 156,
    x: 84,
    y: 195,
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
