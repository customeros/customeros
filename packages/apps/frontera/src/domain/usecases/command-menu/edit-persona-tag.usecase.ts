import { RootStore } from '@store/root';
import { TagDatum, TagStore } from '@store/Tags/Tag.store';
import { action, computed, reaction, observable, runInAction } from 'mobx';

import { Tag, EntityType } from '@graphql/types';

export class EditPersonaTagUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newTags = new Set();
  @observable public accessor initialTags: TagStore[] = [];
  private root = RootStore.getInstance();

  constructor() {
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.computeInitialTags = this.computeInitialTags.bind(this);

    reaction(() => this.newTags.size, this.computeInitialTags);
    reaction(() => this.contextIds, this.computeInitialTags);
    this.computeInitialTags();
  }

  @computed
  get contact() {
    if (this.root.ui.commandMenu.type !== 'EditPersonaTag') return;

    return this.root.contacts.getById(
      this.root.ui.commandMenu.context.ids?.[0] as string,
    );
  }

  @computed
  get contextIds() {
    if (this.root.ui.commandMenu.type !== 'EditPersonaTag') return [];

    return this.root.ui.commandMenu.context.ids;
  }

  @computed
  get inputLabel() {
    const selectedIds = this.contextIds;
    const contact = this.contact;

    return selectedIds?.length === 1
      ? `Contact - ${contact?.value?.name}`
      : `${selectedIds?.length} contacts`;
  }

  @computed
  get contactTags() {
    return new Set((this.contact?.value?.tags ?? []).map((tag) => tag.name));
  }

  @action
  private computeInitialTags() {
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
  }

  @action
  public select(t: TagDatum) {
    if (!this.contact) return;

    if (this.contextIds?.length === 1) {
      const foundIndex = this.contact.value?.tags?.findIndex(
        (e) => e.metadata.id === t.metadata.id,
      );

      this.contact.draft();

      if (typeof foundIndex !== 'undefined' && foundIndex > -1) {
        this.contact.value.tags?.splice(foundIndex, 1);
        this.newTags.delete(t.name);
      } else {
        this.contact.value.tags = this.contact.value.tags ?? [];
        this.contact.value.tags.push(t as TagDatum);
      }
      this.contact.commit();
    } else {
      this.root.contacts.updateTags(this.contextIds, [t as Tag]);
    }
  }

  @action
  public create() {
    const name = this.searchTerm;

    if (!this.contact) return;

    this.root.tags?.create(
      { name, entityType: EntityType.Contact },
      {
        onSucces: (id) => {
          runInAction(() => {
            this.contact?.draft();
            this.contact?.value.tags?.push({
              name,
              metadata: {
                id,
              },
              colorCode:
                this.root.tags.getById(name)?.value?.colorCode ?? 'grayModern',
              entityType: EntityType.Contact,
            });
            this.contact?.commit();

            this.newTags.add(name);
          });

          this.setSearchTerm('');
        },
      },
    );
  }
}
