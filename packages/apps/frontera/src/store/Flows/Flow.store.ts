import type { RootStore } from '@store/root';

import set from 'lodash/set';
import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { Edge, MarkerType } from '@xyflow/react';
import { runInAction, makeAutoObservable } from 'mobx';
import { FlowService } from '@store/Flows/__service__';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';
import { FlowNodeType, FlowActionType } from '@store/Flows/types.ts';

import { Flow, DataSource, FlowStatus } from '@graphql/types';

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
        this.service.mergeFlow({
          input: {
            id: this.id,
            name: this.value.name,
            nodes: JSON.stringify(this.value.nodes),
            edges: JSON.stringify(this.value.edges),
          },
        });
      });
  }

  get parsedNodes() {
    try {
      return JSON.parse(this.value.nodes);
    } catch (error) {
      console.error('Error parsing nodes:', error);

      return initialNodes; // Return an initial array as a fallback
    }
  }

  get parsedEdges() {
    try {
      return JSON.parse(this.value.edges);
    } catch (error) {
      console.error('Error parsing edges:', error);

      return initialEdges; // Return an initial array as a fallback
    }
  }

  public async updateFlow({ nodes, edges }: { nodes: string; edges: string }) {
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
          'Flow updated successfully',
          'update-flow-success',
        );
      });
    } catch (e) {
      runInAction(() => {
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

      return this.root.flowContacts.value.get(item?.metadata?.id)?.value;
    });

    flowContacts && set(output, 'contacts', flowContacts);

    return output;
  }

  public linkContact = async (contactId: string) => {
    this.isLoading = true;

    try {
      const contactStore = this.root.contacts.value.get(contactId);

      if (contactStore?.flow) {
        await contactStore.deleteFlowContact();
      }

      await this.service.addContact({
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
        this.root.ui.toastSuccess(
          `Contact added to '${this.value.name}'`,
          'link-contact-to-flows-success',
        );
        contactStore?.invalidate();
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
        this.root.ui.toastSuccess(
          `${contactIds.length} contacts added to '${this.value.name}'`,
          'link-contacts-to-flows-success',
        );
        this.root.contacts.sync({ action: 'INVALIDATE', ids: contactIds });
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
  contacts: [],
  nodes: JSON.stringify(initialNodes),
  edges: JSON.stringify(initialEdges),
});
const initialNodes = [
  {
    id: 'tn-1',
    type: FlowNodeType.Trigger,
    position: { x: 250, y: 100 },
    data: {
      action: FlowActionType.FLOW_START,
      entity: null,
      triggerType: null,
    },
  },
  {
    id: 'tn-2',
    type: FlowNodeType.Control,
    position: { x: 315, y: 300 },
    data: {
      action: FlowActionType.FLOW_END,
    },
  },
];

const initialEdges: Edge[] = [
  {
    id: 'e1-2',
    source: 'tn-1',
    target: 'tn-2',
    selected: false,
    selectable: true,
    focusable: true,
    interactionWidth: 60,
    markerEnd: {
      type: MarkerType.ArrowClosed,
      width: 60,
      height: 60,
    },
    type: 'baseEdge',
  },
];
