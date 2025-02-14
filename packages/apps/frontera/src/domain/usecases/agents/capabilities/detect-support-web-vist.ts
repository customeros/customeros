import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { CapabilityType } from '@shared/types/__generated__/graphql.types';

export class DetectSupportWebVistUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();

  @observable public accessor supportUrls: string[] = [];
  @observable public accessor pageToTrack: string = '';
  @observable public accessor isOpen: boolean = false;
  @observable public accessor validationError: string = '';

  constructor(private readonly agentId: string) {
    this.open = this.open.bind(this);
    this.close = this.close.bind(this);
    this.toggle = this.toggle.bind(this);
    this.setPageToTrack = this.setPageToTrack.bind(this);
    this.addPageToTrack = this.addPageToTrack.bind(this);
    this.removePageToTrack = this.removePageToTrack.bind(this);

    this.init();
  }

  @computed
  get listenerErrors() {
    return this.root.agents
      .getById(this.agentId)
      ?.value.capabilities.find(
        (c) => c.type === CapabilityType.DetectSupportWebvisit,
      )?.errors;
  }

  @computed
  get isInvalid() {
    return this.pageToTrack.length === 0;
  }

  @action
  init() {
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error('DetectSupportWebVistUsecase.init: Agent not found');

      return;
    }

    const capability = agent.value.capabilities.find(
      (c) => c.type === CapabilityType.DetectSupportWebvisit,
    );

    if (!capability) {
      console.error('DetectSupportWebVistUsecase.init: Capability not found');

      return;
    }

    const config = Agent.parseConfig(capability.config);

    if (!config) {
      console.error('DetectSupportWebVistUsecase.init: Config not found');

      return;
    }

    if (!config.supportUrls) {
      console.error(
        'DetectSupportWebVistUsecase.init: Pages to track not found',
      );

      return;
    }

    this.supportUrls = Array.isArray(config.supportUrls.value)
      ? config.supportUrls.value
      : [];
  }

  @action
  toggle(open: boolean) {
    this.isOpen = open;
  }

  @action
  open() {
    this.isOpen = true;
  }

  @action
  close() {
    this.isOpen = false;
    this.pageToTrack = '';
  }

  @action
  setPageToTrack(page: string) {
    this.pageToTrack = page;
  }

  @action
  addPageToTrack(page: string) {
    const span = Tracer.span('DetectSupportWebVistUsecase.addPageToTrack', {
      page,
      supportUrls: this.supportUrls,
    });

    this.supportUrls.push(page);

    span.end({
      supportUrls: this.supportUrls,
    });
  }

  @action
  removePageToTrack(page: string) {
    const span = Tracer.span('DetectSupportWebVistUsecase.removePageToTrack', {
      page,
      supportUrls: this.supportUrls,
    });

    this.supportUrls = this.supportUrls.filter((p) => p !== page);

    span.end();
  }

  @action
  validate() {
    const span = Tracer.span('DetectSupportWebVistUsecase.validate');

    if (this.pageToTrack.length === 0) {
      console.error('DetectSupportWebVistUsecase.validate: No pages to track');

      this.validationError = 'Huston we have a blank';

      return false;
    }

    span.end();

    return true;
  }

  @action
  async executeAdd(opts?: { onInvalid?: () => void }) {
    const span = Tracer.span('DetectSupportWebVistUsecase.executeAdd');
    const isValid = this.validate();

    if (!isValid) {
      opts?.onInvalid?.();

      span.end();

      return;
    }

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error('DetectSupportWebVistUsecase.executeAdd: Agent not found');

      return;
    }

    this.addPageToTrack(this.pageToTrack);
    this.close();

    span.end();

    agent?.setCapabilityConfig(
      CapabilityType.DetectSupportWebvisit,
      'supportUrls',
      this.supportUrls,
    );

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'DetectSupportWebVistUsecase.executeAdd: Error saving agent',
      );
    }

    if (res) {
      agent.put(res?.agent_Save);
      this.init();
    }

    span.end();
  }

  @action
  async executeRemove(page: string) {
    const span = Tracer.span('DetectSupportWebVistUsecase.executeRemove', {
      page,
    });

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'DetectSupportWebVistUsecase.executeRemove: Agent not found',
      );

      return;
    }

    this.removePageToTrack(page);

    span.end();

    agent?.setCapabilityConfig(
      CapabilityType.DetectSupportWebvisit,
      'supportUrls',
      this.supportUrls,
    );

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'DetectSupportWebVistUsecase.executeRemove: Error saving agent',
      );
    }

    if (res) {
      agent.put(res?.agent_Save);
      this.init();
    }

    span.end();
  }
}
