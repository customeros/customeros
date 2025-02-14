import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { TagStore } from '@store/Tags/Tag.store';
import { action, computed, reaction, observable } from 'mobx';
import { TagService, OrganizationService } from '@domain/services';

import { EntityType } from '@graphql/types';

export class EditOrganizationTagUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newTags = new Set();
  @observable public accessor initialTags: TagStore[] = [];
  @observable public accessor shouldPreventClose = true;

  private root = RootStore.getInstance();
  private tagService = new TagService();
  private organizationService = new OrganizationService();

  constructor() {
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.computeInitialTags = this.computeInitialTags.bind(this);

    reaction(() => this.newTags.size, this.computeInitialTags);
    this.computeInitialTags();
  }

  @computed
  get organization() {
    if (
      !['ChangeTags', 'OrganizationCommands'].includes(
        this.root.ui.commandMenu.type,
      )
    ) {
      return;
    }

    return this.root.organizations.getById(
      this.root.ui.commandMenu.context.ids?.[0] as string,
    );
  }

  @computed
  get contextIds() {
    if (
      !['ChangeTags', 'OrganizationCommands'].includes(
        this.root.ui.commandMenu.type,
      )
    ) {
      return [];
    }

    return this.root.ui.commandMenu.context.ids;
  }

  @computed
  get inputLabel() {
    const selectedIds = this.contextIds;
    const organization = this.organization;

    return selectedIds?.length === 1
      ? `Company - ${organization?.value?.name}`
      : `${selectedIds?.length} companies`;
  }

  @computed
  get organizationTags() {
    return new Set(
      (this.organization?.value?.tags ?? []).map((tag) => tag.name),
    );
  }

  @action
  private computeInitialTags() {
    const span = Tracer.span('EditOrganizationTagUsecase.computeInitialTags');

    this.initialTags = this.root.tags
      ?.getByEntityType(EntityType.Organization)
      .filter((e) => !!e.value.name)
      .sort((a, b) => {
        const aInOrg = this.organizationTags.has(a.value.name);
        const bInOrg = this.organizationTags.has(b.value.name);

        if (aInOrg && !bInOrg) return -1;
        if (!aInOrg && bInOrg) return 1;

        return 0;
      });

    span.end();
  }

  @computed
  get tagList() {
    const sorted = this.initialTags
      .filter((e) => !!e.value.name)
      .sort((a, b) => {
        const aInOrg = this.newTags.has(a.value.name);
        const bInOrg = this.newTags.has(b.value.name);

        if (aInOrg && !bInOrg) return -1;
        if (!aInOrg && bInOrg) return 1;

        return 0;
      });

    return sorted.filter((tag) =>
      tag.value.name.toLowerCase().includes(this.searchTerm.toLowerCase()),
    );
  }

  @action
  public setSearchTerm(searchTerm: string) {
    this.searchTerm = searchTerm;
  }

  @action
  public reset() {
    this.setSearchTerm('');
    this.newTags.clear();
    this.shouldPreventClose = true;
  }

  @action
  public preventClose() {
    this.shouldPreventClose = true;
  }

  @action
  public allowClose() {
    this.shouldPreventClose = false;
  }

  @action
  public close() {
    this.reset();
    this.root.ui.commandMenu.setOpen(false);
  }

  private handleSingleOrganizationTag(tag: TagStore): void {
    const hasTag = this.organization?.value?.tags?.some(
      (t) => t.metadata.id === tag.id,
    );

    if (hasTag) {
      this.organizationService.removeTag(this.organization!, tag);
      this.newTags.delete(tag.value.name);
    } else {
      this.organizationService.addTag(this.organization!, tag);
    }
  }

  private handleMultipleOrganizationsTags(tag: TagStore): void {
    const allOrgsHaveTag = this.contextIds.every((id) => {
      const org = this.root.organizations.getById(id);

      return org?.value.tags?.some((t) => t.metadata.id === tag.id);
    });

    this.contextIds.forEach((id) => {
      const org = this.root.organizations.getById(id);

      if (!org) return;

      allOrgsHaveTag
        ? this.organizationService.removeTag(org, tag)
        : this.organizationService.addTag(org, tag);
    });
  }

  @action
  public select(id?: string) {
    const span = Tracer.span('EditOrganizationTagUsecase.select', { id });

    if (!id || !this.organization) {
      console.error(
        'EditOrganizationTagUsecase: select called without id or company',
      );

      return;
    }

    const tag = this.root.tags.getById(id);

    if (!tag) {
      console.error('EditOrganizationTagUsecase: Tag not found in the store');

      return;
    }

    if (this.contextIds.length === 1) {
      this.handleSingleOrganizationTag(tag);
    } else {
      this.handleMultipleOrganizationsTags(tag);
    }

    if (!this.shouldPreventClose) {
      this.close();
    }

    span.end();
  }

  @action
  public create() {
    const name = this.searchTerm;

    if (!this.organization) return;

    this.tagService.createTag(
      { name, entityType: EntityType.Organization },
      {
        onSuccess: (id) => {
          this.select(id);
          this.newTags.add(name);
          this.setSearchTerm('');

          if (!this.shouldPreventClose) {
            this.close();
          }
        },
      },
    );
  }
}
