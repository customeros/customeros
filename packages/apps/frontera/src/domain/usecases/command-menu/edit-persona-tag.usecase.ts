import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { TagStore } from '@store/Tags/Tag.store';
import { TagService, ContactService } from '@domain/services';
import { action, computed, reaction, observable } from 'mobx';

import { EntityType } from '@graphql/types';

export class EditPersonaTagUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newTags = new Set();
  @observable public accessor initialTags: TagStore[] = [];
  @observable public accessor shouldPreventClose = true;

  private root = RootStore.getInstance();
  private tagService = new TagService();
  private contactService = new ContactService();

  constructor() {
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.computeInitialTags = this.computeInitialTags.bind(this);

    reaction(() => this.newTags.size, this.computeInitialTags);
    this.computeInitialTags();
  }

  @computed
  get contact() {
    if (
      !['EditPersonaTag', 'ContactCommands'].includes(
        this.root.ui.commandMenu.type,
      )
    ) {
      return;
    }

    return this.root.contacts.getById(
      this.root.ui.commandMenu.context.ids?.[0] as string,
    );
  }

  @computed
  get contextIds() {
    if (
      !['EditPersonaTag', 'ContactCommands'].includes(
        this.root.ui.commandMenu.type,
      )
    )
      return [];

    return this.root.ui.commandMenu.context.ids;
  }

  @computed
  get inputLabel() {
    const selectedIds = this.contextIds;
    const contact = this.contact;

    return selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;
  }

  @computed
  get contactTags() {
    return new Set((this.contact?.value?.tags ?? []).map((tag) => tag.name));
  }

  @action
  private computeInitialTags() {
    const span = Tracer.span('EditPersonaTagUsecase.computeInitialTags');

    this.initialTags = this.root.tags
      ?.getByEntityType(EntityType.Contact)
      .filter((e) => !!e.value.name)
      .sort((a, b) => {
        const aInOrg = this.contactTags.has(a.value.name);
        const bInOrg = this.contactTags.has(b.value.name);

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

  private handleSingleContactTag(tag: TagStore): void {
    const hasTag = this.contact?.value?.tags?.some(
      (t) => t.metadata.id === tag.id,
    );

    if (hasTag) {
      this.contactService.removeTag(this.contact!, tag);
      this.newTags.delete(tag.value.name);
    } else {
      this.contactService.addTag(this.contact!, tag);
    }
  }

  private handleMultipleContactsTags(tag: TagStore): void {
    const allContactsHaveTag = this.contextIds.every((id) => {
      const contact = this.root.contacts.getById(id);

      return contact?.value.tags?.some((t) => t.metadata.id === tag.id);
    });

    this.contextIds.forEach((id) => {
      const contact = this.root.contacts.getById(id);

      if (!contact) {
        console.error('EditContactTagUsecase: Contact not found in the store');

        return;
      }

      allContactsHaveTag
        ? this.contactService.removeTag(contact, tag)
        : this.contactService.addTag(contact, tag);
    });
  }

  @action
  public select(id?: string) {
    if (!id || !this.contact) {
      console.error(
        'EditContactTagUsecase: called select without id or contact',
      );

      return;
    }

    const tag = this.root.tags.getById(id);

    if (!tag) {
      console.error('EditContactTagUsecase: Tag not found in the store');

      return;
    }

    if (this.contextIds.length === 1) {
      this.handleSingleContactTag(tag);
    } else {
      this.handleMultipleContactsTags(tag);
    }

    if (!this.shouldPreventClose) {
      this.close();
    }
  }

  @action
  public create() {
    const name = this.searchTerm;

    if (!this.contact) return;

    this.tagService.createTag(
      { name, entityType: EntityType.Contact },
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
