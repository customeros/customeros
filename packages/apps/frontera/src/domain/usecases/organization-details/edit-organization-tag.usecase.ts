import { RootStore } from '@store/root';
import { action, computed, reaction, observable } from 'mobx';
import { TagService, OrganizationService } from '@domain/services';

import { EntityType } from '@graphql/types';
import { SelectOption } from '@ui/utils/types.ts';

export class EditOrganizationTagUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newTags = new Set();
  @observable public accessor initialTags: { label: string; value: string }[] =
    [];

  private root = RootStore.getInstance();
  private tagService = new TagService();
  private organizationService = new OrganizationService();
  private organizationId: string;

  constructor(organizationId: string) {
    this.organizationId = organizationId;
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.computeInitialTags = this.computeInitialTags.bind(this);

    reaction(
      () =>
        this.newTags.size ||
        this.root.tags.getByEntityType(EntityType.Organization).length,

      this.computeInitialTags,
    );
    this.computeInitialTags();
  }

  @computed
  get organization() {
    if (!this.organizationId) return;

    return this.root.organizations.getById(this.organizationId);
  }

  @computed
  get orgTags() {
    return new Set(
      (this.organization?.value?.tags ?? []).map((tag) => tag.name),
    );
  }

  @computed
  get selectedTags() {
    return (
      this.organization?.value?.tags.map(
        (tag) =>
          ({
            label: tag.name,
            value: tag?.metadata.id,
          } as SelectOption),
      ) ?? []
    );
  }

  @action
  private computeInitialTags() {
    this.initialTags = this.root.tags
      .getByEntityType(EntityType.Organization)
      .filter((e) => !!e.value.name)
      .sort((a, b) => {
        const aInOrg = this.orgTags.has(a.value.name);
        const bInOrg = this.orgTags.has(b.value.name);

        if (aInOrg && !bInOrg) return -1;
        if (!aInOrg && bInOrg) return 1;

        return 0;
      })
      .map((tag) => {
        return {
          value: tag.id,
          label: tag.tagName,
        };
      });
  }

  @computed
  get tagList() {
    const sorted = this.initialTags.slice().sort((a, b) => {
      const aInOrg = this.newTags.has(a.label);
      const bInOrg = this.newTags.has(b.label);

      if (aInOrg && !bInOrg) return -1;
      if (!aInOrg && bInOrg) return 1;

      return 0;
    });

    return sorted.filter((tag) =>
      tag.label.toLowerCase().includes(this.searchTerm.toLowerCase()),
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
  }

  @action
  public close() {
    this.reset();
    this.root.ui.commandMenu.setOpen(false);
  }

  @action
  public select(ids?: string[]) {
    const organization = this.organization;

    if (!ids || !organization) {
      console.error(
        'EditOrganizationTagUsecase: select called without id or company',
      );

      return;
    }

    const tags = ids.map((e) => this.root.tags.getById(e));
    const currentTags = organization.value.tags;

    if (!tags || !tags.length) {
      currentTags.forEach(({ metadata }) => {
        const tag = this.root.tags.getById(metadata.id);

        if (!tag) {
          console.error(
            'EditOrganizationTagUsecase: Tag not found in the store',
          );

          return;
        }
        this.organizationService.removeTag(organization, tag);
      });
      this.newTags.clear();

      return;
    }

    const tagsToAdd = tags
      .filter(
        (tag) =>
          !currentTags.some((currentTag) => currentTag.metadata.id === tag?.id),
      )
      .filter(Boolean);

    const tagsToRemove = currentTags
      .filter(
        (currentTag) => !tags.some((tag) => tag?.id === currentTag.metadata.id),
      )
      .map((tagDatum) => this.root.tags.getById(tagDatum.metadata.id))
      .filter(Boolean);

    tagsToAdd.forEach((tag) => {
      if (!tag) {
        console.error('EditOrganizationTagUsecase: Tag not found in the store');

        return;
      }
      this.organizationService.addTag(organization, tag);
    });

    tagsToRemove.forEach((tag) => {
      if (!tag) {
        console.error('EditOrganizationTagUsecase: Tag not found in the store');

        return;
      }
      this.newTags.delete(tag.value.name);
      this.organizationService.removeTag(organization, tag);
    });
  }

  @action
  public create() {
    const name = this.searchTerm;

    if (!this.organization) return;

    this.tagService.createTag(
      { name, entityType: EntityType.Organization },
      {
        onSuccess: (id) => {
          this.select([id]);
          this.newTags.add(name);
          this.setSearchTerm('');
        },
      },
    );
  }
}
