import { MasterPlanMilestone } from '@graphql/types';

export type MilestoneDatum = Omit<
  MasterPlanMilestone,
  'appSource' | 'source' | 'sourceOfTruth' | 'createdAt' | 'updatedAt'
>;
