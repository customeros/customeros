import { runInAction } from 'mobx';
import { RootStore } from '@store/root';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service';

import {
  SortingDirection,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export class OrganizationService {
  private root = RootStore.getInstance();
  private service = OrganizationsService.getInstance();

  constructor() {}

  public async searchTenant(searchTerm: string) {
    try {
      const { ui_organizations_search } =
        await this.service.searchOrganizations({
          limit: 30,
          sort: {
            by: 'ORGANIZATIONS_NAME',
            caseSensitive: false,
            direction: SortingDirection.Asc,
          },
          where: {
            OR: [
              {
                filter: {
                  property: 'ORGANIZATIONS_NAME',
                  value: searchTerm,
                  caseSensitive: false,
                  includeEmpty: false,
                  operation: ComparisonOperator.Contains,
                },
              },
              {
                filter: {
                  property: 'ORGANIZATIONS_WEBSITE',
                  value: searchTerm,
                  caseSensitive: false,
                  includeEmpty: false,
                  operation: ComparisonOperator.Contains,
                },
              },
            ],
          },
        });

      const results = ui_organizations_search?.ids ?? [];

      if (results.length === 0) {
        return;
      }

      await this.root.organizations.retrieve(results);
    } catch (_err) {
      throw new Error('Failed to search for tenant');
    }
  }

  public async merge(primaryId: string, secondaryIds: string[]) {
    try {
      const { organization_Merge } = await this.service.mergeOrganizations({
        primaryOrganizationId: primaryId,
        mergedOrganizationIds: secondaryIds,
      });

      runInAction(() => {
        if (organization_Merge.id) {
          this.root.organizations.mergeOrganizations(primaryId, secondaryIds);

          this.root.ui.toastSuccess(
            `Merged organizations`,
            `merge-${primaryId}`,
          );
        }
      });
    } catch (err) {
      throw new Error('Failed to merge organizations');
    }
  }
}
