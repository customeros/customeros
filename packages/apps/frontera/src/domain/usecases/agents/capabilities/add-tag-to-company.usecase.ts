import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { TagService } from '@domain/services';
import { Agent } from '@store/Agents/Agent.dto';
import { action, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';

import { CapabilityType } from '@graphql/types';
import { EntityType } from '@shared/types/__generated__/graphql.types';

export class AddTagToCompanyUsecase {
  @observable public accessor searchTerm = '';
  @observable public accessor newTags = new Set();
  @observable public accessor initialTags: { label: string; value: string }[] =
    [];

  private root = RootStore.getInstance();
  private tagService = new TagService();
  private agentService = new AgentService();

  constructor(private agentId: string) {
    this.agentId = agentId;
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);

    this.init();
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
          this.select(id);
          this.newTags.add(name);
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

    this.initialTags = this.tagList
      .filter((t) => t.label === parsedCapability.tagName.value)
      .map((t) => ({ label: t.label, value: t.value }));

    span.end({
      initialTags: this.initialTags,
    });
  }

  @action
  public select(id?: string) {
    const span = Tracer.span('AddTagToCompanyUsecase.select', {
      id,
    });

    if (!id) {
      this.reset();

      return;
    }

    this.newTags.clear();
    this.newTags.add(id);

    span.end({
      ids: Array.from(this.newTags),
    });
  }

  @computed
  get tagList() {
    return this.root.tags
      .getByEntityType(EntityType.Organization)
      .filter((e) => !!e.value.name)
      .map((tag) => ({
        label: tag.tagName,
        value: tag.id,
      }));
  }

  @computed
  get selectedTags() {
    return this.tagList.filter(
      (tag) =>
        this.newTags.has(tag.value) ||
        this.initialTags.some((t) => t.value === tag.value),
    );
  }

  @action
  public reset() {
    this.searchTerm = '';
    this.newTags.clear();
    this.initialTags = [];
  }

  @action
  public async execute() {
    const span = Tracer.span('AddTagToCompanyUsecase.execute');

    const agent = this.root.agents.getById(this.agentId);

    if (!agent) {
      console.error(
        'AddTagToCompanyUsecase.execute: Agent not found. aborting execution',
      );

      return;
    }

    const tagName = this.selectedTags.map((tag) => tag.label).join(', ');

    agent?.setCapabilityConfig(
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

    if (res) {
      agent.put(res?.agent_Save);
      this.init();
    }

    span.end();
  }
}
