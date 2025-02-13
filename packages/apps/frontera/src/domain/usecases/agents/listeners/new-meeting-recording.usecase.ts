import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { AgentListenerEvent } from '@graphql/types';

export class NewMeetingRecordingUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();
  private hasInitted = false;

  @observable accessor isOpen: boolean = false;
  @observable accessor webhookUrl: string = '';
  @observable accessor listenerErrors: string = '';
  @observable private accessor _notetaker: string = '';

  constructor(private agentId: string) {
    this.setNotetaker = this.setNotetaker.bind(this);
    this.init();
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
  toggle() {
    const span = Tracer.span('NewMeetingRecordingUsecase.toggle', {
      previous: this.isOpen,
    });

    this.isOpen = !this.isOpen;

    span.end({
      current: this.isOpen,
    });
  }

  init() {
    if (this.hasInitted) {
      return;
    }

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

    const listenerConfig = Agent.parseConfig<'webhookUrl' | 'meetingSource'>(
      listener.config,
    );

    if (!listenerConfig) {
      console.error(
        'NewMeetingRecordingUsecase.init: Listener config could not be parsed, aborting.',
      );

      span.end();

      return;
    }

    // if (!listenerConfig.webhookUrl) {
    //   console.error(
    //     'NewMeetingRecordingUsecase.init: Webhook URL not found, aborting.',
    //   );

    //   span.end();

    //   return;
    // }

    if (!listenerConfig.meetingSource) {
      console.error(
        'NewMeetingRecordingUsecase.init: Meeting source not found, aborting.',
      );

      span.end();

      return;
    }

    // this.setWebhookUrl(listenerConfig.webhookUrl.value as string);
    this.setNotetaker(listenerConfig.meetingSource.value as string);

    this.hasInitted = true;
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

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'NewMeetingRecordingUsecase.execute: Error updating agent, aborting.',
      );
    }

    if (res) {
      agent.put(res.agent_Save);
      this.init();
    }

    span.end();
  }

  static notetakerOptions = [
    { label: 'Grain', value: 'grain' },
    { label: 'Fathom', value: 'fathom' },
  ];
}
