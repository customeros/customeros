import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { validateUrl } from '@utils/url';
import { CapabilityType } from '@graphql/types';

export class AddWebsiteToTrackUsecase {
  private service = new AgentService();
  private root = RootStore.getInstance();

  @observable accessor website: string = '';
  @observable accessor isOpen: boolean = false;
  @observable accessor websites: string[] = [];
  @observable accessor validationError: string = '';

  constructor(private agentId: string) {
    this.toggle = this.toggle.bind(this);
    this.open = this.open.bind(this);
    this.close = this.close.bind(this);
    this.execute = this.execute.bind(this);
    this.validate = this.validate.bind(this);
    this.setWebsite = this.setWebsite.bind(this);
    this.removeWebsite = this.removeWebsite.bind(this);
  }

  @computed
  get isInvalid() {
    return this.validationError.length > 0;
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
    this.websites.push(this.website);
  }

  @action
  removeWebsite(website: string) {
    this.websites = this.websites.filter((w) => w !== website);
  }

  @action
  validate() {
    if (this.website.length === 0) {
      this.validationError = 'Houston we have a blank';

      return false;
    }

    if (!validateUrl(this.website)) {
      this.validationError = 'This domain appears to be invalid';

      return false;
    }

    if (this.websites.includes(this.website)) {
      this.validationError = 'This website is already added';

      return false;
    }

    return true;
  }

  execute(opts?: { onInvalid?: () => void }) {
    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error('AddWebsiteToTrack: Agent not found. aborting execution');

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

      agent?.setCapabilityConfig(
        CapabilityType.IdentifyWebVisitor,
        'websites',
        this.websites,
      );

      this.service.saveAgent(agent);
    }
  }
}
