import { match } from 'ts-pattern';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';

import { ColumnViewType } from '@graphql/types';

export const getOpportunitiesSortFn = (columnId: string) =>
  match(columnId)
    .with(ColumnViewType.OpportunitiesName, () => (row: OpportunityStore) => {
      return row.value?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.OpportunitiesOrganization,
      () => (row: OpportunityStore) =>
        row.organization?.value?.name?.trim().toLocaleLowerCase() || null,
    )
    .with(
      ColumnViewType.OpportunitiesStage,
      () => (row: OpportunityStore) => row.value.externalStage || null,
    )
    .with(ColumnViewType.OpportunitiesOwner, () => (row: OpportunityStore) => {
      return row.owner?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.OpportunitiesTimeInStage,
      () => (row: OpportunityStore) => {
        return row.value.stageLastUpdated
          ? new Date(row.value.stageLastUpdated)
          : null;
      },
    )
    .with(
      ColumnViewType.OpportunitiesCreatedDate,
      () => (row: OpportunityStore) => {
        return row.value.metadata.created
          ? new Date(row.value.metadata.created)
          : null;
      },
    )
    .with(
      ColumnViewType.OpportunitiesEstimatedArr,
      () => (row: OpportunityStore) => {
        return row.value.maxAmount;
      },
    )
    .with(
      ColumnViewType.OpportunitiesNextStep,
      () => (row: OpportunityStore) => {
        return row.value.nextSteps;
      },
    )
    .otherwise(() => (_row: OpportunityStore) => false);
