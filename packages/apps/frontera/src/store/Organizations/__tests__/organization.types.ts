// types/organization.types.ts

export interface OnboardingDetails {
  status: string;
  updatedAt: null;
  comments: string;
}

export interface RenewalSummary {
  arrForecast: null;
  maxArrForecast: null;
  nextRenewalDate: null;
  renewalLikelihood: null;
}

export interface AccountDetails {
  ltv: number;
  churned: null;
  onboarding: OnboardingDetails;
  renewalSummary: RenewalSummary;
}

export interface LastTouchpoint {
  lastTouchPointAt: Date | null;
  lastTouchPointType: string | null;
  lastTouchPointTimelineEvent: string | null;
  lastTouchPointTimelineEventId: string | null;
}

export interface ParentCompany {
  id: string;
  name: string;
  relationship?: string;
}

export interface SocialMedia {
  url: string;
}

export interface Tag {
  name: string;
}

export interface Subsidiary {
  organization: {
    name: string;
  };
}

export interface Organization {
  owner: null;
  icon: string;
  logo: string;
  name: string;
  stage: string;
  contracts: null;
  public: boolean;
  website: string;
  industry: string;
  domains: string[];
  employees: number;
  yearFounded: null;
  leadSource: string;
  tags: Tag[] | null;
  description: string;
  isCustomer: boolean;
  locations: string[];
  relationship: string;
  valueProposition: string;
  socialMedia: SocialMedia[];
  subsidiaries: Subsidiary[];
  accountDetails: AccountDetails;
  lastTouchpoint: LastTouchpoint;
  parentCompanies: ParentCompany[];
}

export type AssertionType = 'not.toBeNull' | 'toBeNull';

export interface SpecialAssertion {
  assertType: AssertionType;
}

export type ExpectedValue =
  | string
  | number
  | boolean
  | null
  | SpecialAssertion
  | Record<string, unknown>
  | unknown[];

export interface ExpectedState {
  [key: string]: ExpectedValue;
}
