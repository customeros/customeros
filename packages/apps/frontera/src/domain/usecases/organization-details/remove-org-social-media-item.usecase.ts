import { action } from 'mobx';
import { RootStore } from '@store/root';
import { OrganizationService } from '@domain/services';

import { Social } from '@graphql/types';

export class RemoveOrgSocialMediaItemUsecase {
  private root = RootStore.getInstance();
  private service = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.remove = this.remove.bind(this);
    this.orgId = orgId;
  }

  @action
  public async remove(id: string) {
    const organization = this.root.organizations.getById(this.orgId);

    if (!id || !organization) {
      console.error(
        'RemoveOrgSocialMediaItemUsecase: remove social called without id or company',
      );

      return;
    }
    const social =
      (organization?.value?.socialMedia?.find((s) => s.id === id) as Social) ??
      null;

    if (!social) {
      console.error(
        'RemoveOrgSocialMediaItemUsecase: social not found on company',
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
