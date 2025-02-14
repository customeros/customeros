import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { WebhookIntegration } from '@infra/repositories/common';
import { AgentService } from '@domain/services/agent/agent.service';
import { CommonService } from '@domain/services/common/common.service';

import { AgentListenerEvent } from '@graphql/types';

export class NewMeetingRecordingUsecase {
  private service = new AgentService();
  private commonService = new CommonService();
  private root = RootStore.getInstance();
  private hasInitted = false;

  @observable accessor isOpen: boolean = false;
  @observable accessor webhookUrl: string = '';
  @observable accessor listenerErrors: string = '';
  @observable private accessor _notetaker: string = '';
  @observable private accessor agentId: string = '';
  @observable accessor isWebhookUrlLoading: boolean = false;

  constructor() {
    this.setNotetaker = this.setNotetaker.bind(this);
  }

  @computed
  get notetaker() {
    return NewMeetingRecordingUsecase.notetakerOptions.find(
      (option) => option.value === this._notetaker,
    );
  }

  @action
  setNotetaker(notetaker: string) {
    const span = Tracer.span('NewMeetingRecordingUsecase.setNotetaker', {
      previous: this._notetaker,
    });

    this._notetaker = notetaker;

    span.end({
      current: notetaker,
    });
  }

  @action
  setWebhookUrl(webhookUrl: string) {
    const span = Tracer.span('NewMeetingRecordingUsecase.setWebhookUrl', {
      previous: this.webhookUrl,
    });

    this.webhookUrl = webhookUrl;

    span.end({
      current: webhookUrl,
    });
  }

  @action
  setAgentId(agentId: string) {
    const span = Tracer.span('NewMeetingRecordingUsecase.setAgentId', {
      previous: this.agentId,
    });

    this.agentId = agentId;

    span.end({
      current: agentId,
    });
  }

  @action
  toggle() {
    const span = Tracer.span('NewMeetingRecordingUsecase.toggle', {
      previous: this.isOpen,
    });

    this.isOpen = !this.isOpen;

    span.end({
      current: this.isOpen,
    });
  }

  @action
  setIsWebhookUrlLoading(isWebhookUrlLoading: boolean) {
    this.isWebhookUrlLoading = isWebhookUrlLoading;
  }

  @action
  reset() {
    this.hasInitted = false;
    this.agentId = '';
    this.webhookUrl = '';
    this.listenerErrors = '';
    this._notetaker = '';
    this.isWebhookUrlLoading = false;
  }

  private async initWebhookUrl(integration: WebhookIntegration) {
    const span = Tracer.span('NewMeetingRecordingUsecase.initWebhookUrl', {
      integration,
    });

    this.setIsWebhookUrlLoading(true);

    const [res, err] = await this.commonService.getWebhookUrl(integration);

    if (err) {
      console.error(
        'NewMeetingRecordingUsecase.initWebhookUrl: Error getting webhook URL, aborting.',
        err,
      );

      this.setIsWebhookUrlLoading(false);
      span.end();

      return;
    }

    if (res) {
      this.setWebhookUrl(res.data.hook.url);
    }

    this.setIsWebhookUrlLoading(false);

    span.end();
  }

  private async createWebhookUrl(integration: WebhookIntegration) {
    const span = Tracer.span('NewMeetingRecordingUsecase.createWebhookUrl', {
      integration,
    });

    this.setIsWebhookUrlLoading(true);

    const [webhookRes, webhookErr] = await this.commonService.createWebhookUrl(
      this.notetaker?.value as WebhookIntegration,
    );

    if (webhookErr) {
      console.error(
        'NewMeetingRecordingUsecase.createWebhookUrl: Error creating webhook URL, aborting.',
        webhookErr,
      );

      span.end();

      this.setIsWebhookUrlLoading(false);

      return;
    }

    if (webhookRes) {
      this.setWebhookUrl(webhookRes.data.hook.url);
    }

    this.setIsWebhookUrlLoading(false);

    span.end();
  }

  async init(agentId: string) {
    if (this.hasInitted) {
      return;
    }

    this.hasInitted = true;
    this.setAgentId(agentId);

    const span = Tracer.span('NewMeetingRecordingUsecase.init', {
      agentId: this.agentId,
    });

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'NewMeetingRecordingUsecase.init: Agent not found, aborting.',
      );
      span.end();

      return;
    }

    const listener = agent.value.listeners.find(
      (l) => l.type === AgentListenerEvent.NewMeetingRecording,
    );

    if (!listener) {
      console.error(
        'NewMeetingRecordingUsecase.init: Listener not found, aborting.',
      );

      span.end();

      return;
    }

    const listenerConfig = Agent.parseConfig<'meetingSource'>(listener.config);

    if (!listenerConfig) {
      console.error(
        'NewMeetingRecordingUsecase.init: Listener config could not be parsed, aborting.',
      );

      span.end();

      return;
    }

    if (!listenerConfig.meetingSource) {
      console.error(
        'NewMeetingRecordingUsecase.init: Meeting source not found, aborting.',
      );

      span.end();

      return;
    }

    if (listenerConfig.meetingSource.value) {
      this.setNotetaker(listenerConfig.meetingSource.value as string);

      await this.initWebhookUrl(
        listenerConfig.meetingSource.value as WebhookIntegration,
      );
    }

    span.end();
  }

  async execute() {
    const span = Tracer.span('NewMeetingRecordingUsecase.execute', {
      agentId: this.agentId,
    });

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'NewMeetingRecordingUsecase.execute: Agent not found, aborting.',
      );
      span.end();

      return;
    }

    agent.setListenerConfig(
      AgentListenerEvent.NewMeetingRecording,
      'meetingSource',
      this.notetaker?.value,
    );

    await this.createWebhookUrl(this.notetaker?.value as WebhookIntegration);

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'NewMeetingRecordingUsecase.execute: Error updating agent, aborting.',
      );
    }

    if (res) {
      agent.put(res.agent_Save);
      this.init(this.agentId);
    }

    span.end();
  }

  static notetakerOptions = [
    { label: 'Grain', value: 'grain' },
    { label: 'Fathom', value: 'fathom' },
  ];
}
