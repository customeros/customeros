import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { CapabilityType } from '@graphql/types';

export const cooldownPeriods = [1, 4, 8, 12, 24, 999999];
export const cooldownPeriodsMap: Record<number, string> = {
  1: '1h',
  4: '4h',
  8: '8h',
  12: '12h',
  24: '24h',
  999999: `Don't notify again`,
};

export class AddSlackChannelUsecase {
  private root = RootStore.getInstance();
  private agentService = new AgentService();

  @observable accessor isOpen = false;
  @observable accessor inputValue = '';
  @observable accessor selectedChannel = '';
  @observable accessor cooldownPeriod = 12;
  @observable accessor slackChannelError = '';

  constructor(private agentId: string) {
    this.togglePopover = this.togglePopover.bind(this);
    this.close = this.close.bind(this);
    this.selectChannel = this.selectChannel.bind(this);
    this.setInputValue = this.setInputValue.bind(this);
    this.enableSlack = this.enableSlack.bind(this);
    this.disableSlack = this.disableSlack.bind(this);
    this.toggleSendSlackNotificationActive =
      this.toggleSendSlackNotificationActive.bind(this);
    this.execute = this.execute.bind(this);

    this.init();
  }

  @computed
  get allOptions() {
    return this.root.common.slackChannelsArray.map((val) => {
      return { value: val.channelId, label: val.name };
    });
  }

  @computed
  get options() {
    return this.allOptions
      .filter((option) =>
        option?.label.toLowerCase().includes(this.inputValue.toLowerCase()),
      )
      .map((option) => ({
        value: option.value,
        label: option.label,
      }));
  }

  @computed
  get selectedOption() {
    return this.allOptions.find(
      (option) => option.value === this.selectedChannel,
    );
  }

  @computed
  get selectedChannelName() {
    return this.selectedOption?.label;
  }

  @computed
  get isSlackEnabled() {
    return this.root.settings.slack.enabled;
  }

  @computed
  get slackNotificationActive() {
    return this.root.agents
      .getById(this.agentId)
      ?.value.capabilities?.find(
        (c) => c.type === CapabilityType.WebVisitorSendSlackNotification,
      )?.active;
  }

  @computed
  get capabilityErrors() {
    return this.root.agents
      .getById(this.agentId)
      ?.value.capabilities.find(
        (c) => c.type === CapabilityType.IdentifyWebVisitor,
      )?.errors;
  }

  @action
  togglePopover(open: boolean) {
    this.isOpen = open;
  }

  @action
  close() {
    const span = Tracer.span('AddSlackChannelUsecase.close', {
      isOpen: this.isOpen,
    });

    this.isOpen = false;

    span.end(`isOpen = ${this.isOpen}`);
  }

  @action
  selectChannel(channel: string) {
    const span = Tracer.span('AddSlackChannelUsecase.selectChannel', {
      previous: this.selectedChannel,
    });

    this.selectedChannel = channel;
    this.close();
    this.execute();

    span.end({ selectedChannel: this.selectedChannel });
  }

  @action
  setInputValue(value: string) {
    this.inputValue = value;
  }

  @action
  setCooldownPeriod(period: number) {
    const span = Tracer.span('AddSlackChannelUsecase.setCooldownPeriod', {
      period,
    });

    this.cooldownPeriod = period;
    this.execute();

    span.end(`this.cooldownPeriod = ${this.cooldownPeriod}`);
  }

  @action
  init() {
    const span = Tracer.span('AddSlackChannelUsecase.init', {
      selectedChannel: this.selectedChannel,
      cooldownPeriod: this.cooldownPeriod,
    });
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error('AddSlackChannelUsecase.setDefaultValues: Agent not found');

      return;
    }

    const foundCapabilityConfig = agent.value?.capabilities?.find(
      (c) => c.type === CapabilityType.WebVisitorSendSlackNotification,
    )?.config;

    if (!foundCapabilityConfig) {
      console.error(
        'AddSlackChannelUsecase.setDefaultValues: Capability not found',
      );

      return;
    }

    const parsedCapability = Agent.parseConfig(foundCapabilityConfig);

    if (!parsedCapability) {
      console.error(
        'AddSlackChannelUsecase.setDefaultValues: Could not parse capability',
      );

      return;
    }

    if (
      !Object.hasOwn(parsedCapability, 'channelId') ||
      !Object.hasOwn(parsedCapability, 'cooldownHours')
    ) {
      console.error(
        'AddSlackChannelUsecase.setDefaultValues: Missing default properties',
      );

      return;
    }

    this.selectedChannel = parsedCapability.channelId.value as string;
    this.cooldownPeriod = parsedCapability.cooldownHours.value as number;
    this.slackChannelError = parsedCapability.channelId.error as string;
    span.end({
      selectedChannel: this.selectedChannel,
      cooldownPeriod: this.cooldownPeriod,
      slackChannelError: this.slackChannelError,
    });
  }

  enableSlack() {
    const capabilityId = new URLSearchParams(window.location.search).get('cid');

    if (!capabilityId) {
      console.error(
        'AddSlackChannelUsecase.enableSlack: Capability ID not found',
      );

      return;
    }

    this.root.settings.slack.enableSync(
      `${import.meta.env.VITE_CLIENT_APP_URL}/agents`,
      encodeURIComponent(`${this.agentId}:cid=${capabilityId}`),
    );
  }

  disableSlack() {
    this.root.settings.slack.disableSync();
  }

  async toggleSendSlackNotificationActive() {
    const span = Tracer.span(
      'AddSlackChannelUsecase.toggleSendSlackNotificationActive',
    );

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddSlackChannelUsecase.toggleSendSlackNotificationActive: Agent not found',
      );

      return;
    }

    agent?.toggleCapabilityStatus(
      CapabilityType.WebVisitorSendSlackNotification,
    );

    const [res, err] = await this.agentService.saveAgent(agent);

    if (err) {
      agent?.toggleCapabilityStatus(
        CapabilityType.WebVisitorSendSlackNotification,
      );

      console.error(
        'AddSlackChannelUsecase.toggleSendSlackNotificationActive: Error saving agent. aborting execution',
      );
    }

    if (res) {
      agent.put(res?.agent_Save);
      this.init();
    }

    span.end();
  }

  async execute() {
    const span = Tracer.span('AddSlackChannelUsecase.execute');

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error('AddSlackChannelUsecase.execute: Agent not found');

      return;
    }

    agent?.setCapabilityConfig(
      CapabilityType.WebVisitorSendSlackNotification,
      'channelId',
      this.selectedChannel,
    );

    agent?.setCapabilityConfig(
      CapabilityType.WebVisitorSendSlackNotification,
      'cooldownHours',
      this.cooldownPeriod,
    );

    const [res, err] = await this.agentService.saveAgent(agent);

    if (err) {
      console.error(
        'AddSlackChannelUsecase.execute: Error saving agent. aborting execution',
      );
    }

    if (res) {
      agent.put(res?.agent_Save);
      this.init();
    }

    span.end();
  }
}
