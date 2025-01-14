import { RootStore } from '@store/root';
import { action, observable, runInAction } from 'mobx';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service';

export class RemoveOrgSocialMediaItemUsecase {
  @observable accessor organizationId: string = '';
  private root = RootStore.getInstance();
  private service = OrganizationsService.getInstance();

  constructor() {
    this.remove = this.remove.bind(this);
    this.setId = this.setId.bind(this);
  }

  @action
  setId(id: string) {
    this.organizationId = id;
  }

  @action
  public async remove(id: string) {
    try {
      await this.service.removeSocial({ socialId: id });

      runInAction(() => {
        const org = this.root.organizations.getById(this.organizationId);
        const idx = org?.value?.socialMedia.findIndex((s) => s.id === id);

        if (typeof idx === 'undefined' || idx < 0) return;
        org?.value?.socialMedia?.splice(idx, 1);
      });
    } catch {
      this.root.ui.toastError(
        "We couldn't remove this link",
        `${id}-remove-social-media`,
      );
    } finally {
      this.root.organizations.sync({
        action: 'INVALIDATE',
        ids: [this.organizationId],
      });
      this.organizationId = '';
    }
  }
}
