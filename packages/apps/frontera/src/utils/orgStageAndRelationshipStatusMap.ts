import { OrganizationStage, OrganizationRelationship } from '@graphql/types';

export const stageRelationshipMap: Record<
  OrganizationStage,
  OrganizationRelationship
> = {
  [OrganizationStage.Unqualified]: OrganizationRelationship.NotAFit,
  [OrganizationStage.Lead]: OrganizationRelationship.Prospect,
  [OrganizationStage.Target]: OrganizationRelationship.Prospect, // Default to PROSPECT for simplicity
  [OrganizationStage.Engaged]: OrganizationRelationship.Prospect,
  [OrganizationStage.Trial]: OrganizationRelationship.Prospect,
  [OrganizationStage.ReadyToBuy]: OrganizationRelationship.Prospect,
  [OrganizationStage.Onboarding]: OrganizationRelationship.Customer,
  [OrganizationStage.InitialValue]: OrganizationRelationship.Customer,
  [OrganizationStage.RecurringValue]: OrganizationRelationship.Customer,
  [OrganizationStage.MaxValue]: OrganizationRelationship.Customer,
  [OrganizationStage.PendingChurn]: OrganizationRelationship.Customer,
};

export const relationshipStageMap: Record<
  OrganizationRelationship,
  OrganizationStage
> = {
  [OrganizationRelationship.NotAFit]: OrganizationStage.Unqualified,
  [OrganizationRelationship.Prospect]: OrganizationStage.Target,
  [OrganizationRelationship.Customer]: OrganizationStage.Onboarding,
  [OrganizationRelationship.FormerCustomer]: OrganizationStage.Target,
};

export const validRelationshipsForStage: Record<
  OrganizationStage,
  Array<OrganizationRelationship>
> = {
  [OrganizationStage.Unqualified]: [OrganizationRelationship.NotAFit],
  [OrganizationStage.Lead]: [OrganizationRelationship.Prospect],
  [OrganizationStage.Target]: [
    OrganizationRelationship.Prospect,
    OrganizationRelationship.FormerCustomer,
  ],
  [OrganizationStage.Engaged]: [OrganizationRelationship.Prospect],
  [OrganizationStage.Trial]: [OrganizationRelationship.Prospect],
  [OrganizationStage.ReadyToBuy]: [OrganizationRelationship.Prospect],
  [OrganizationStage.Onboarding]: [OrganizationRelationship.Customer],
  [OrganizationStage.InitialValue]: [OrganizationRelationship.Customer],
  [OrganizationStage.RecurringValue]: [OrganizationRelationship.Customer],
  [OrganizationStage.MaxValue]: [OrganizationRelationship.Customer],
  [OrganizationStage.PendingChurn]: [OrganizationRelationship.Customer],
};
