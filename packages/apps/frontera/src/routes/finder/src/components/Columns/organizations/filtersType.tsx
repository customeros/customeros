import { Cake } from '@ui/media/icons/Cake';
import { Key01 } from '@ui/media/icons/Key01';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Hash02 } from '@ui/media/icons/Hash02';
import { Globe01 } from '@ui/media/icons/Globe01';
import { Users03 } from '@ui/media/icons/Users03';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Activity } from '@ui/media/icons/Activity';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Building07 } from '@ui/media/icons/Building07';
import { Calculator } from '@ui/media/icons/Calculator';
import { Building05 } from '@ui/media/icons/Building05';
import { ArrowCircleDownRight } from '@ui/media/icons/ArrowCircleDownRight';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import {
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';
export type FilterType = {
  // options?: any[];
  icon: JSX.Element;
  filterName: string;
  filterAccesor: ColumnViewType;
  filterOperators: ComparisonOperator[];
  filterType: 'text' | 'date' | 'number' | 'list';
};

// import { type RootStore } from '@store/root';

// const getFilterTypes=(store:RootStore)=>{

// }

export const filterTypes: Partial<Record<ColumnViewType, FilterType>> = {
  [ColumnViewType.OrganizationsName]: {
    filterType: 'text',
    filterAccesor: ColumnViewType.OrganizationsName,
    filterName: 'Organization name',
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Building07 />,
    // options: test.organizations.toArray(),
  },
  [ColumnViewType.OrganizationsWebsite]: {
    filterType: 'text',
    filterAccesor: ColumnViewType.OrganizationsWebsite,
    filterName: 'Website',
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Globe01 />,
  },
  [ColumnViewType.OrganizationsRelationship]: {
    filterType: 'list',
    filterName: 'Relationship',
    filterAccesor: ColumnViewType.OrganizationsRelationship,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <AlignHorizontalCentre02 />,
  },
  [ColumnViewType.OrganizationsRenewalLikelihood]: {
    filterType: 'list',
    filterName: 'Health',
    filterAccesor: ColumnViewType.OrganizationsRenewalLikelihood,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Activity />,
  },
  [ColumnViewType.OrganizationsOnboardingStatus]: {
    filterType: 'text',
    filterName: 'Onboarding status',
    filterAccesor: ColumnViewType.OrganizationsOnboardingStatus,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Trophy01 />,
  },
  [ColumnViewType.OrganizationsRenewalDate]: {
    filterType: 'date',
    filterName: 'Renewal Date',
    filterAccesor: ColumnViewType.OrganizationsRenewalDate,
    filterOperators: [ComparisonOperator.Gt, ComparisonOperator.Lt],
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
    filterType: 'list',
    filterName: 'Owner',
    filterAccesor: ColumnViewType.OrganizationsOwner,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Key01 />,
  },
  [ColumnViewType.OrganizationsLeadSource]: {
    filterType: 'list',
    filterName: 'Source',
    filterAccesor: ColumnViewType.OrganizationsLeadSource,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <ArrowCircleDownRight />,
  },
  [ColumnViewType.OrganizationsCreatedDate]: {
    filterType: 'date',
    filterName: 'Created Date',
    filterAccesor: ColumnViewType.OrganizationsCreatedDate,
    filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
    icon: <Calendar />,
  },
  [ColumnViewType.OrganizationsYearFounded]: {
    filterType: 'number',
    filterName: 'Founded',
    filterAccesor: ColumnViewType.OrganizationsYearFounded,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.Eq,
      ComparisonOperator.NotEqual,
    ],
    icon: <Cake />,
  },
  [ColumnViewType.OrganizationsEmployeeCount]: {
    filterType: 'number',
    filterName: 'Employees',
    filterAccesor: ColumnViewType.OrganizationsEmployeeCount,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.Eq,
      ComparisonOperator.NotEqual,
    ],
    icon: <Users03 />,
  },
  [ColumnViewType.OrganizationsSocials]: {
    filterType: 'number',
    filterName: 'Employees',
    filterAccesor: ColumnViewType.OrganizationsSocials,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.Eq,
      ComparisonOperator.NotEqual,
    ],
    icon: <Calendar />,
  },
  [ColumnViewType.OrganizationsLastTouchpoint]: {
    filterType: 'date',
    filterName: 'Last touchpoint',
    filterAccesor: ColumnViewType.OrganizationsLastTouchpoint,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Calendar />,
  },
  [ColumnViewType.OrganizationsChurnDate]: {
    filterType: 'date',
    filterName: 'Churn date',
    filterAccesor: ColumnViewType.OrganizationsChurnDate,
    filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
    icon: <Calendar />,
  },
  [ColumnViewType.OrganizationsLastTouchpointDate]: {
    filterType: 'date',
    filterName: 'Last interacted',
    filterAccesor: ColumnViewType.OrganizationsLastTouchpointDate,
    filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
    icon: <Calendar />,
  },
  [ColumnViewType.OrganizationsLtv]: {
    filterType: 'number',
    filterName: 'LTV',
    filterAccesor: ColumnViewType.OrganizationsLtv,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.NotEqual,
      ComparisonOperator.Eq,
    ],
    icon: <CurrencyDollarCircle />,
  },
  [ColumnViewType.OrganizationsIndustry]: {
    filterType: 'list',
    filterName: 'Industry',
    filterAccesor: ColumnViewType.OrganizationsIndustry,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Building05 />,
  },
  [ColumnViewType.OrganizationsContactCount]: {
    filterType: 'number',
    filterName: 'Contact count',
    filterAccesor: ColumnViewType.OrganizationsContactCount,
    filterOperators: [
      ComparisonOperator.Lt,
      ComparisonOperator.Gt,
      ComparisonOperator.NotEqual,
      ComparisonOperator.Eq,
    ],
    icon: <Hash02 />,
  },
  [ColumnViewType.OrganizationsTags]: {
    filterType: 'list',
    filterName: 'Tags',
    filterAccesor: ColumnViewType.OrganizationsTags,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Tag01 />,
  },
  [ColumnViewType.OrganizationsIsPublic]: {
    filterType: 'list',
    filterName: 'Ownership type',
    filterAccesor: ColumnViewType.OrganizationsIsPublic,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Key01 />,
  },
  [ColumnViewType.OrganizationsStage]: {
    filterType: 'list',
    filterName: 'Stage',
    filterAccesor: ColumnViewType.OrganizationsStage,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Columns03 />,
  },
  [ColumnViewType.OrganizationsCity]: {
    filterType: 'list',
    filterName: 'Country',
    filterAccesor: ColumnViewType.OrganizationsCity,
    filterOperators: [
      ComparisonOperator.Contains,
      ComparisonOperator.IsEmpty,
      ComparisonOperator.NotContains,
      ComparisonOperator.IsNotEmpty,
    ],
    icon: <Building07 />,
  },
};
