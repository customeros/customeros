import { RootStore } from '@store/root.ts';
import { TagService, ContactService } from '@domain/services';
import { action, computed, reaction, observable } from 'mobx';

import { EntityType } from '@graphql/types';
import { SelectOption } from '@ui/utils/types.ts';

export class EditContactTagUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newTags = new Set();
  @observable public accessor initialTags: SelectOption[] = [];

  private root = RootStore.getInstance();
  private tagService = new TagService();
  private contactService = new ContactService();

  private contactId: string;

  constructor(contactId: string) {
    this.contactId = contactId;
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.computeInitialTags = this.computeInitialTags.bind(this);

    reaction(
      () =>
        this.newTags.size ||
        this.root.tags.getByEntityType(EntityType.Contact).length,
      this.computeInitialTags,
    );
    this.computeInitialTags();
  }

  @computed
  get contact() {
    if (!this.contactId) return;

    const contact = this.root.contacts.getById(this.contactId);

    if (!contact) {
      return;
    }

    return contact;
  }

  @computed
  get contactTags() {
    return new Set((this.contact?.value?.tags ?? []).map((tag) => tag.name));
  }

  @computed
  get selectedTags() {
    return (
      this.contact?.value?.tags.map(
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
    const tags = this.root.tags?.getByEntityType(EntityType.Contact) ?? [];

    this.initialTags = tags
      .sort((a, b) => {
        const aInOrg = this.contactTags.has(a.value.name);
        const bInOrg = this.contactTags.has(b.value.name);

        if (aInOrg && !bInOrg) return -1;
        if (!aInOrg && bInOrg) return 1;

        return 0;
      })
      .map(
        (tag) =>
          ({
            label: tag.value.name,
            value: tag.id,
          } as SelectOption),
      );
  }

  @computed
  get tagList() {
    const sorted = this.initialTags
      .filter((e) => !!e.label)
      .sort((a, b) => {
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
  }

  @action
  public select(ids?: string[]) {
    const contact = this.contact;

    if (!ids || !contact) {
      console.error('EditContactTagUsecase: No contact or tag ids provided');

      return;
    }

    const currentTags = this.contact.value.tags;
    const tags = ids.map((e) => this.root.tags.getById(e)).filter(Boolean);

    if (!tags.length) {
      currentTags.forEach(({ metadata }) => {
        const tag = this.root.tags.getById(metadata?.id);

        if (!tag) {
          console.error('EditContactTagUsecase: Tag not found in the store');

          return;
        }

        this.contactService.removeTag(contact, tag);
      });
      this.newTags.clear();

      return;
    }

    const selectedTagNames = new Set(tags.map((tag) => tag?.value.name));

    currentTags.forEach(({ metadata }) => {
      const tag = this.root.tags.getById(metadata?.id);

      if (!tag) {
        console.error('EditContactTagUsecase: Tag not found in the store');

        return;
      }

      if (!selectedTagNames.has(tag.value.name)) {
        this.contactService.removeTag(contact, tag);
        this.newTags.delete(tag.value.name);
      }
    });

    tags.forEach((tag) => {
      if (!tag) {
        console.error('EditContactTagUsecase: Tag not found in the store');

        return;
      }
      const foundIndex = contact.value.tags.findIndex(
        (e) => e.metadata.id === tag?.id,
      );

      if (foundIndex === -1) {
        this.contactService.addTag(contact, tag);
      }
    });
  }

  @action
  public create() {
    const name = this.searchTerm.trim();

    if (!name || !this.contact) return;

    this.tagService.createTag(
      { name, entityType: EntityType.Contact },
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
