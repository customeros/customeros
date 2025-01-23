import { RootStore } from '@store/root';
import { action, computed } from 'mobx';
import { OrganizationService } from '@domain/services';

import { Social } from '@graphql/types';
import { getOrganizationUUID } from '@utils/getOrganizationUUID';

export class RemoveOrgSocialMediaItemUsecase {
  private root = RootStore.getInstance();
  private service = new OrganizationService();

  constructor() {
    this.remove = this.remove.bind(this);
  }

  @computed
  get organization() {
    const uuid = getOrganizationUUID(window.location.pathname);

    if (!uuid) {
      console.error('Invalid usage of RemoveOrgSocialMediaItemUsecase');

      return;
    }

    return this.root.organizations.getById(uuid);
  }

  @action
  public async remove(id: string) {
    const organization = this.organization;

    if (!id || !organization) {
      console.error(
        'RemoveOrgSocialMediaItemUsecase: remove social called without id or organization',
      );

      return;
    }
    const social =
      (organization?.value?.socialMedia?.find((s) => s.id === id) as Social) ??
      null;

    if (!social) {
      console.error(
        'RemoveOrgSocialMediaItemUsecase: social not found on organization',
      );

      return;
    }

    const [res] = await this.service.removeSocialMediaItem(
      organization,
      social,
    );

    if (res) {
      this.root.organizations.sync({
        action: 'INVALIDATE',
        ids: [organization.id],
      });
    }
  }
}
