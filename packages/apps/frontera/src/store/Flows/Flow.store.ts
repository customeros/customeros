import type { RootStore } from '@store/root';

import set from 'lodash/set';
import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { MarkerType } from '@xyflow/react';
import { Transport } from '@infra/transport';
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

  async startFlow(options?: { onSuccess?: () => void }) {
    this.isLoading = true;

    try {
      const { flow_On } = await this.service.startFlow({
        id: this.id,
      });

      runInAction(() => {
        if (flow_On?.metadata?.id) {
          this.value.status = FlowStatus.On;
          options?.onSuccess && options?.onSuccess();
        }
      });
    } catch (error) {
      this.root.ui.toastError("We couldn't start the flow", 'start-flow-error');
    } finally {
      runInAction(() => {
        this.invalidate();
      });
      this.isLoading = false;
    }
  }

  async stopFlow({ onSuccess }: { onSuccess?: () => void }) {
    this.isLoading = true;

    try {
      const { flow_Off } = await this.service.stopFlow({
        id: this.id,
      });

      runInAction(() => {
        if (flow_Off?.metadata?.id) {
          this.value.status = FlowStatus.Off;
          onSuccess && onSuccess();
        }
      });
    } catch (error) {
      this.root.ui.toastError("We couldn't stop the flow", 'stop-flow-error');
    } finally {
      runInAction(() => {
        this.invalidate();
      });
      this.isLoading = false;
    }
  }

  async updateFlowName() {
    this.isLoading = true;

    try {
      await this.service.changeFlowName({
        id: this.id,
        name: this.value.name,
      });
    } catch (error) {
      this.root.ui.toastError(
        "We couldn't update flow name",
        'update-flow-name-error',
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

    match(path).with(['name', ...P.array()], () => {
      this.updateFlowName();
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
          `Flow changes published`,
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
          "We couldn't publish the flow changes",
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
        contactStore?.draft();
        contactStore?.value.flows.push(this.value.metadata.id);
        contactStore?.commit({ syncOnly: true });

        const newFLowContact = new FlowParticipantStore(
          this.root,
          this.transport,
        );

        const newFlowContactValue = {
          ...newFLowContact.value,
          metadata: {
            ...newFLowContact.value.metadata,
            id: flowParticipant_Add.metadata.id,
          },
          entityType: 'CONTACT',
          entityId: contactId,
          status: FlowParticipantStatus.Scheduled,
          executions: [],
          requirementsUnmeet: [],
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
        contactStores.forEach((contactStore) => {
          contactStore?.draft();
          contactStore!.value.flows = [
            ...(contactStore?.value.flows ?? []),
            this.value.name,
          ];
          contactStore?.commit({ syncOnly: true });
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

  public sendTestEmail = async (
    payload: {
      subject: string;
      bodyTemplate: string;
      sendToEmailAddress: string;
    },
    options?: {
      onSuccess?: () => void;
      onFinally?: () => void;
    },
  ) => {
    this.isLoading = true;

    try {
      await this.service.sendTestEmail(payload);

      this.root.ui.toastSuccess('Test email sent', 'send-test-email-success');
      runInAction(() => {
        options?.onSuccess?.();
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          "We couldn't send the test email",
          'send-test-email-error',
        );
      });
    } finally {
      runInAction(() => {
        options?.onFinally?.();
        this.isLoading = false;
      });
    }
  };
}

const getDefaultValue = (): Flow => ({
  name: '',
  status: FlowStatus.Off,
  description: '',
  tableViewDefId: '',
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
