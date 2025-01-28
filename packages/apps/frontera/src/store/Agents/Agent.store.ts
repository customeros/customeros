import { Store } from '@store/_store';
import { type RootStore } from '@store/root';
import { type Transport } from '@infra/transport';
import { type AgentDatum } from '@infra/repositories/agent';

import { CapabilityType } from '@graphql/types';

import { Agent } from './Agent.dto';

export class AgentStore extends Store<AgentDatum, Agent> {
  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Agents',
      getId: (data) => data.id,
      factory: Agent,
    });

    const agent = new Agent(
      this,
      Agent.default({
        name: 'Web visit identifier',
        goal: 'Identify website visitors and add them as enriched leads to CustomerOS',
        capabilities: [
          {
            id: '1',
            name: 'Track and identify website visitors',
            action: '',
            type: CapabilityType.WebVisitorSendSlackNotification,
            // values: JSON.stringify({
            //   websites: {
            //     value: [],
            //   },
            // }),
            errors: null,
          },
          {
            id: '2',
            name: 'Log page views and session duration',
            action: '',
            type: CapabilityType.WebVisitorSendSlackNotification,
            // values: '',
            errors: null,
          },
          {
            id: '3',
            name: 'Create identified organizations as leads',
            action: '',
            type: CapabilityType.WebVisitorSendSlackNotification,
            // values: '',
            errors: null,
          },
          {
            id: '4',
            name: 'Analyze behaviour for intent signals',
            action: '',
            type: CapabilityType.WebVisitorSendSlackNotification,
            // values: '',
            errors: null,
          },
          {
            id: '5',
            name: 'Send Slack notification',
            action: '',
            type: CapabilityType.SendSlackNotification,
            // values: JSON.stringify({
            //   slackId: {
            //     value: [],
            //   },
            // }),
            errors: null,
          },
        ],
      }),
    );

    this.value.set(agent.id, agent);
  }
}
