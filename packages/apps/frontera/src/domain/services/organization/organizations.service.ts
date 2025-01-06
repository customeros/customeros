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
      console.info(_err);
    }
  }
}
