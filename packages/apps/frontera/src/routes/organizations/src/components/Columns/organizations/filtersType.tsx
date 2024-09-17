import { Key01 } from '@ui/media/icons/Key01';
import { Globe01 } from '@ui/media/icons/Globe01';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Building07 } from '@ui/media/icons/Building07';
import { Calculator } from '@ui/media/icons/Calculator';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import {
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export const filterTypes = {
  [ColumnViewType.OrganizationsName]: {
    filterType: 'text',
    filterAccesor: ColumnViewType.OrganizationsName,
    filterName: 'Organization name',
    filterOperators: [ComparisonOperator.Contains, ComparisonOperator.IsEmpty],
    icon: <Building07 />,
  },
  [ColumnViewType.OrganizationsWebsite]: {
    filterType: 'text',
    filterAccesor: ColumnViewType.OrganizationsWebsite,
    filterName: 'Website',
    filterOperators: [ComparisonOperator.Contains, ComparisonOperator.IsEmpty],
    icon: <Globe01 />,
  },
  [ColumnViewType.OrganizationsRelationship]: {
    filterType: 'text',
    filterName: 'Relationship',
    filterAccesor: ColumnViewType.OrganizationsRelationship,
    filterOperators: [ComparisonOperator.Eq, ComparisonOperator.IsEmpty],
    icon: <AlignHorizontalCentre02 />,
  },
  [ColumnViewType.OrganizationsOnboardingStatus]: {
    filterType: 'text',
    filterName: 'Onboarding status',
    filterAccesor: ColumnViewType.OrganizationsOnboardingStatus,
    filterOperators: [ComparisonOperator.Eq, ComparisonOperator.IsEmpty],
    icon: <Trophy01 />,
  },
  [ColumnViewType.OrganizationsRenewalDate]: {
    filterType: 'date',
    filterName: 'Renewal Date',
    filterAccesor: ColumnViewType.OrganizationsRenewalDate,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.Between,
    ],
    icon: <Calendar />,
  },
  [ColumnViewType.OrganizationsForecastArr]: {
    filterType: 'number',
    filterName: 'Arr forecast',
    filterAccesor: ColumnViewType.OrganizationsForecastArr,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.Between,
      ComparisonOperator.Eq,
    ],
    icon: <Calculator />,
  },
  [ColumnViewType.OrganizationsOwner]: {
    filterType: 'text',
    filterName: 'Owner',
    filterAccesor: ColumnViewType.OrganizationsOwner,
    filterOperators: [ComparisonOperator.Eq, ComparisonOperator.IsEmpty],
    icon: <Key01 />,
  },
  [ColumnViewType.OrganizationsLeadSource]: {},
};
