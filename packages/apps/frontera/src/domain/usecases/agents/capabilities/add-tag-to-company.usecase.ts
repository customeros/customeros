import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { TagService } from '@domain/services';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { CapabilityType } from '@graphql/types';
import { EntityType } from '@shared/types/__generated__/graphql.types';

type Tag = { label: string; value: string };

export class AddTagToCompanyUsecase {
  @observable public accessor searchTerm = '';
  @observable public accessor newTags = new Set<string>();
  @observable public accessor initialTag: Tag | null = null;

  private root = RootStore.getInstance();
  private tagService = new TagService();
  private agentService = new AgentService();

  constructor(private readonly agentId: string) {
    this.create = this.create.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.execute = this.execute.bind(this);

    this.init();
  }

  @computed
  get listenerErrors() {
    return this.root.agents
      .getById(this.agentId)
      ?.value.capabilities.find(
        (c) => c.type === CapabilityType.ApplyTagToCompany,
      )?.errors;
  }

  @action
  setSearchTerm(searchTerm: string) {
    this.searchTerm = searchTerm;
  }

  @action
  public create() {
    const name = this.searchTerm;

    if (!this.agentId) return;

    this.tagService.createTag(
      { name, entityType: EntityType.Organization },
      {
        onSuccess: (id) => {
          this.newTags.add(id);
          this.execute({ value: id, label: '' });
          this.setSearchTerm('');
        },
      },
    );
  }

  @action
  init() {
    const span = Tracer.span('AddTagToCompanyUsecase.init', {
      searchTerm: this.searchTerm,
    });

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddTagToCompanyUsecase.init: Agent not found. aborting execution',
      );

      return;
    }

    const foundCapabilityConfig = agent.value?.capabilities?.find(
      (c) => c.type === CapabilityType.ApplyTagToCompany,
    )?.config;

    if (!foundCapabilityConfig) {
      console.error(
        'AddTagToCompanyUsecase.init: Capability not found. aborting execution',
      );

      return;
    }

    const parsedCapability = Agent.parseConfig(foundCapabilityConfig);

    if (!parsedCapability) {
      console.error(
        'AddTagToCompanyUsecase.init: Could not parse capability. aborting execution',
      );

      return;
    }

    if (!Object.hasOwn(parsedCapability, 'tagName')) {
      console.error(
        'AddTagToCompanyUsecase.init: Missing tags property. aborting execution',
      );

      return;
    }

    const matchingTag = this.tagList.find(
      (t) => t.label === parsedCapability.tagName.value,
    );

    this.initialTag = matchingTag || null;

    span.end({
      initialTag: this.initialTag,
    });
  }

  @computed
  get tagList(): Tag[] {
    return this.root.tags
      .getByEntityType(EntityType.Organization)
      .filter((e) => !!e.value.name)
      .map((tag) => ({
        label: tag.tagName,
        value: tag.id,
      }))
      .filter((tag) =>
        tag.label.toLowerCase().includes(this.searchTerm.toLowerCase()),
      );
  }

  @computed
  get selectedTag(): Tag | null {
    const id = this.initialTag?.value;

    if (!id) return null;

    const tag = this.root.tags.getById(id);

    if (!tag) return null;

    return {
      label: tag.tagName,
      value: id,
    };
  }

  @action
  public reset() {
    this.searchTerm = '';
    this.newTags.clear();
    this.initialTag = null;
  }

  @action
  public async execute(option?: Tag) {
    const span = Tracer.span('AddTagToCompanyUsecase.execute');

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddTagToCompanyUsecase.execute: Agent not found. aborting execution',
      );

      return;
    }

    if (!option?.value) {
      this.reset();
    }

    if (option?.value) {
      this.newTags.clear();
      this.newTags.add(option.value);
      this.initialTag = option;
    }

    const tagName = this.selectedTag?.label || '';

    agent.setCapabilityConfig(
      CapabilityType.ApplyTagToCompany,
      'tagName',
      tagName,
    );

    const [res, err] = await this.agentService.saveAgent(agent);

    if (err) {
      console.error(
        'AddTagToCompanyUsecase.execute: Error saving agent. aborting execution',
      );
    }

    if (res?.agent_Save) {
      agent.put(res.agent_Save);
      this.init();
    }

    span.end();
  }
}
