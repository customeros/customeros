import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { validateUrl } from '@utils/url';
import { AgentListenerEvent } from '@graphql/types';

export class AddWebsiteToTrackUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();

  @observable accessor website: string = '';
  @observable accessor isOpen: boolean = false;
  @observable accessor websites: string[] = [];
  @observable accessor validationError: string = '';
  @observable accessor websitesError: string = '';

  constructor(private agentId: string) {
    this.toggle = this.toggle.bind(this);
    this.open = this.open.bind(this);
    this.close = this.close.bind(this);
    this.executeAdd = this.executeAdd.bind(this);
    this.executeRemove = this.executeRemove.bind(this);
    this.validate = this.validate.bind(this);
    this.setWebsite = this.setWebsite.bind(this);
    this.removeWebsite = this.removeWebsite.bind(this);

    this.init();
  }

  @computed
  get isInvalid() {
    return this.validationError.length > 0;
  }

  @computed
  get listenerErrors() {
    return this.root.agents
      .getById(this.agentId)
      ?.value.listeners.find((c) => c.type === AgentListenerEvent.NewWebSession)
      ?.errors;
  }

  @action
  setWebsite(website: string) {
    this.website = website;
  }

  @action
  open() {
    this.isOpen = true;
  }

  @action
  toggle(open: boolean) {
    this.isOpen = open;
  }

  @action
  close() {
    this.isOpen = false;
    this.website = '';
    this.validationError = '';
  }

  @action
  addWebsite() {
    const span = Tracer.span('AddWebsiteToTrackUsecase.addWebsite', {
      websites: this.websites,
    });

    this.websites.push(this.website);

    span.end({
      websites: this.websites,
    });
  }

  @action
  removeWebsite(website: string) {
    const span = Tracer.span('AddWebsiteToTrackUsecase.removeWebsite', {
      websites: this.websites,
    });

    this.websites = this.websites.filter((w) => w !== website);

    span.end({
      websites: this.websites,
    });
  }

  @action
  validate() {
    const span = Tracer.span('AddWebsiteToTrackUsecase.validate');

    if (this.website.length === 0) {
      this.validationError = 'Houston we have a blank';
      span.end();

      return false;
    }

    if (!validateUrl(this.website)) {
      this.validationError = 'This domain appears to be invalid';
      span.end();

      return false;
    }

    if (this.websites.includes(this.website)) {
      this.validationError = 'This website is already added';
      span.end();

      return false;
    }
    span.end();

    return true;
  }

  @action
  init() {
    const span = Tracer.span('AddWebsiteToTrackUsecase.init', {
      websites: this.websites,
    });
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddWebsiteToTrackUsecase.init: Agent not found. aborting execution',
      );

      return;
    }

    const listener = agent.value.listeners.find(
      (c) => c.type === AgentListenerEvent.NewWebSession,
    );

    if (!listener) {
      console.error(
        'AddWebsiteToTrackUsecase.init: Listener not found. aborting execution',
      );

      return;
    }

    const config = Agent.parseConfig(listener.config);

    if (!config) {
      console.error(
        'AddWebsiteToTrackUsecase.init: Could not parse config. aborting',
      );

      return;
    }

    if (!config.websites) {
      console.error(
        'AddWebsiteToTrackUsecase.init: Websites not found in config. aborting execution',
      );

      return;
    }

    this.websites = config.websites.value as string[];
    this.websitesError = config.websites.error ?? '';
    span.end({
      websites: config.websites.value,
    });
  }

  async executeRemove(website: string) {
    const span = Tracer.span('AddWebsiteToTrackUsecase.executeRemove', {
      websites: this.websites,
    });

    this.removeWebsite(website);

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddWebsiteToTrackUsecase.executeRemove: Agent not found. aborting execution',
      );

      return;
    }

    agent?.setListenerConfig(
      AgentListenerEvent.NewWebSession,
      'websites',
      this.websites,
    );

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'AddWebsiteToTrackUsecase.executeRemove: Error saving agent. aborting execution',
      );
    }

    if (res) {
      agent.put(res?.agent_Save);
      this.init();
    }

    span.end({
      websites: this.websites,
    });
  }

  async executeAdd(opts?: { onInvalid?: () => void }) {
    const span = Tracer.span('AddWebsiteToTrackUsecase.executeAdd', {
      websites: this.websites,
    });
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddWebsiteToTrackUsecase.executeAdd: Agent not found. aborting execution',
      );

      return;
    }

    const isValid = this.validate();

    if (!isValid) {
      opts?.onInvalid?.();

      return;
    }

    if (isValid) {
      this.addWebsite();
      this.close();

      agent?.setListenerConfig(
        AgentListenerEvent.NewWebSession,
        'websites',
        this.websites,
      );

      const [res, err] = await this.service.saveAgent(agent);

      if (err) {
        console.error(
          'AddWebsiteToTrackUsecase.executeAdd: Error saving agent. aborting execution',
        );
      }

      if (res) {
        agent.put(res?.agent_Save);
        this.init();
      }

      span.end({
        websites: this.websites,
      });
    }

    span.end();
  }
}
