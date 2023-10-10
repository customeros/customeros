import { GroupedOption, SelectOption } from '@shared/types/SelectOptions';
import { FundingRound } from '@graphql/types';

export const relationshipOptions: SelectOption<boolean>[] = [
  {
    label: 'Customer',
    value: true,
  },
  {
    label: 'Prospect',
    value: false,
  },
];

export const industryOptions: GroupedOption[] = [
  {
    label: 'Energy',
    options: [
      {
        label: 'Energy Equipment & Services',
        value: 'Energy Equipment & Services',
      },
      {
        label: 'Oil, Gas & Consumable Fuels',
        value: 'Oil, Gas & Consumable Fuels',
      },
    ],
  },
  {
    label: 'Materials',
    options: [
      {
        label: 'Chemicals',
        value: 'Chemicals',
      },
      {
        label: 'Construction Materials',
        value: 'Construction Materials',
      },
      {
        label: 'Containers & Packaging',
        value: 'Containers & Packaging',
      },
      {
        label: 'Metals & Mining',
        value: 'Metals & Mining',
      },
      {
        label: 'Paper & Forest Products',
        value: 'Paper & Forest Products',
      },
    ],
  },
  {
    label: 'Industrials',
    options: [
      {
        label: 'Aerospace & Defense',
        value: 'Aerospace & Defense',
      },
      {
        label: 'Building Products',
        value: 'Building Products',
      },
      {
        label: 'Construction & Engineering',
        value: 'Construction & Engineering',
      },
      {
        label: 'Electrical Equipment',
        value: 'Electrical Equipment',
      },
      {
        label: 'Industrial Conglomerates',
        value: 'Industrial Conglomerates',
      },
      {
        label: 'Machinery',
        value: 'Machinery',
      },
      {
        label: 'Trading Companies & Distributors',
        value: 'Trading Companies & Distributors',
      },
      {
        label: 'Commercial Services & Supplies',
        value: 'Commercial Services & Supplies',
      },
      {
        label: 'Professional Services',
        value: 'Professional Services',
      },
      {
        label: 'Air Freight & Logistics',
        value: 'Air Freight & Logistics',
      },
      {
        label: 'Passenger Airlines',
        value: 'Passenger Airlines',
      },
      {
        label: 'Marine Transportation',
        value: 'Marine Transportation',
      },
      {
        label: 'Ground Transportation',
        value: 'Ground Transportation',
      },
      {
        label: 'Transportation Infrastructure',
        value: 'Transportation Infrastructure',
      },
    ],
  },
  {
    label: 'Consumer Discretionary',
    options: [
      {
        label: 'Automobile Components',
        value: 'Automobile Components',
      },
      {
        label: 'Automobiles',
        value: 'Automobiles',
      },
      {
        label: 'Household Durables',
        value: 'Household Durables',
      },
      {
        label: 'Leisure Products',
        value: 'Leisure Products',
      },
      {
        label: 'Textiles, Apparel & Luxury Goods',
        value: 'Textiles, Apparel & Luxury Goods',
      },
      {
        label: 'Hotels, Restaurants & Leisure',
        value: 'Hotels, Restaurants & Leisure',
      },
      {
        label: 'Diversified Consumer Services',
        value: 'Diversified Consumer Services',
      },
      {
        label: 'Distributors',
        value: 'Distributors',
      },
      {
        label: 'Broadline Retail',
        value: 'Broadline Retail',
      },
      {
        label: 'Specialty Retail',
        value: 'Specialty Retail',
      },
    ],
  },
  {
    label: 'Consumer Staples',
    options: [
      {
        label: 'Consumer Staples Distribution & Retail',
        value: 'Consumer Staples Distribution & Retail',
      },
      {
        label: 'Beverages',
        value: 'Beverages',
      },
      {
        label: 'Food Products',
        value: 'Food Products',
      },
      {
        label: 'Tobacco',
        value: 'Tobacco',
      },
      {
        label: 'Household Products',
        value: 'Household Products',
      },
      {
        label: 'Personal Products',
        value: 'Personal Products',
      },
    ],
  },
  {
    label: 'Health Care',
    options: [
      {
        label: 'Health Care Equipment & Supplies',
        value: 'Health Care Equipment & Supplies',
      },
      {
        label: 'Health Care Providers & Services',
        value: 'Health Care Providers & Services',
      },
      {
        label: 'Health Care Technology',
        value: 'Health Care Technology',
      },
      {
        label: 'Biotechnology',
        value: 'Biotechnology',
      },
      {
        label: 'Pharmaceuticals',
        value: 'Pharmaceuticals',
      },
      {
        label: 'Life Sciences Tools & Services',
        value: 'Life Sciences Tools & Services',
      },
    ],
  },
  {
    label: 'Financials',
    options: [
      {
        label: 'Banks',
        value: 'Banks',
      },
      {
        label: 'Financial Services',
        value: 'Financial Services',
      },
      {
        label: 'Consumer Finance',
        value: 'Consumer Finance',
      },
      {
        label: 'Capital Markets',
        value: 'Capital Markets',
      },
      {
        label: 'Mortgage Real Estate Investment Trusts (REITs)',
        value: 'Mortgage Real Estate Investment Trusts (REITs)',
      },
      {
        label: 'Insurance',
        value: 'Insurance',
      },
    ],
  },
  {
    label: 'Information Technology',
    options: [
      {
        label: 'Internet Software & Services',
        value: 'Internet Software & Services',
      },
      {
        label: 'IT Services',
        value: 'IT Services',
      },
      {
        label: 'Software',
        value: 'Software',
      },
      {
        label: 'Communications Equipment',
        value: 'Communications Equipment',
      },
      {
        label: 'Technology Hardware, Storage & Peripherals',
        value: 'Technology Hardware, Storage & Peripherals',
      },
      {
        label: 'Electronic Equipment, Instruments & Components',
        value: 'Electronic Equipment, Instruments & Components',
      },
      {
        label: 'Semiconductors & Semiconductor Equipment',
        value: 'Semiconductors & Semiconductor Equipment',
      },
    ],
  },
  {
    label: 'Communication Services',
    options: [
      {
        label: 'Diversified Telecommunication Services',
        value: 'Diversified Telecommunication Services',
      },
      {
        label: 'Wireless Telecommunication Services',
        value: 'Wireless Telecommunication Services',
      },
      {
        label: 'Media',
        value: 'Media',
      },
      {
        label: 'Entertainment',
        value: 'Entertainment',
      },
      {
        label: 'Interactive Media & Services',
        value: 'Interactive Media & Services',
      },
    ],
  },
  {
    label: 'Utilities',
    options: [
      {
        label: 'Electric Utilities',
        value: 'Electric Utilities',
      },
      {
        label: 'Gas Utilities',
        value: 'Gas Utilities',
      },
      {
        label: 'Multi-Utilities',
        value: 'Multi-Utilities',
      },
      {
        label: 'Water Utilities',
        value: 'Water Utilities',
      },
      {
        label: 'Independent Power and Renewable Electricity Producers',
        value: 'Independent Power and Renewable Electricity Producers',
      },
    ],
  },
  {
    label: 'Real Estate',
    options: [
      {
        label: 'Diversified REITs',
        value: 'Diversified REITs',
      },
      {
        label: 'Industrial REITs',
        value: 'Industrial REITs',
      },
      {
        label: 'Hotel & Resort REITs',
        value: 'Hotel & Resort REITs',
      },
      {
        label: 'Office REITs',
        value: 'Office REITs',
      },
      {
        label: 'Health Care REITs',
        value: 'Health Care REITs',
      },
      {
        label: 'Residential REITs',
        value: 'Residential REITs',
      },
      {
        label: 'Retail REITs',
        value: 'Retail REITs',
      },
      {
        label: 'Specialized REITs',
        value: 'Specialized REITs',
      },
      {
        label: 'Real Estate Management & Development',
        value: 'Real Estate Management & Development',
      },
    ],
  },
];

export const employeesOptions: SelectOption<number>[] = [
  { label: '1 - 20 employees', value: 20 },
  { label: '21 - 50 employees', value: 50 },
  { label: '51 - 100 employees', value: 100 },
  { label: '101 - 250 employees', value: 250 },
  { label: '251 - 500 employees', value: 500 },
  { label: '501 - 2500 employees', value: 2500 },
  { label: '2501 - 5000 employees', value: 5000 },
  { label: '5001 - 10000 employees', value: 10000 },
  { label: '10000+ employees', value: 10001 },
];

export const businessTypeOptions: SelectOption[] = [
  { label: 'B2B', value: 'B2B' },
  { label: 'B2C', value: 'B2C' },
  { label: 'Marketplace', value: 'MARKETPLACE' },
];

export const lastFundingRoundOptions: SelectOption<FundingRound>[] = [
  { label: 'Pre-Seed', value: FundingRound.PreSeed },
  { label: 'Seed', value: FundingRound.Seed },
  { label: 'Series A', value: FundingRound.SeriesA },
  { label: 'Series B', value: FundingRound.SeriesB },
  { label: 'Series C', value: FundingRound.SeriesC },
  { label: 'Series D', value: FundingRound.SeriesD },
  { label: 'Series E', value: FundingRound.SeriesE },
  { label: 'Series F', value: FundingRound.SeriesF },
  { label: 'IPO', value: FundingRound.Ipo },
  { label: 'Friends and Family', value: FundingRound.FriendsAndFamily },
  { label: 'Angel', value: FundingRound.Angel },
  { label: 'Bridge', value: FundingRound.Bridge },
];
