// __tests__/Organizations.test.ts

import { it, expect, describe } from 'vitest';

import { Transport } from '../../transport';
import { OrganizationsService } from '../__service__/Organizations.service';
import {
  Organization,
  ExpectedState,
  ExpectedValue,
  SpecialAssertion,
} from './organization.types';

const transport = new Transport();
const service = OrganizationsService.getInstance(transport);

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

function createFieldPath(parts: string[]): string {
  return parts.join('.');
}

function makeAssertion(actual: unknown, expected: unknown, fieldPath: string) {
  if (Array.isArray(expected)) {
    expect
      .soft(Array.isArray(actual), `Field '${fieldPath}' should be an array`)
      .toBe(true);

    if (Array.isArray(actual)) {
      expect
        .soft(actual.length, `Field '${fieldPath}' length mismatch`)
        .toEqual(expected.length);

      expected.forEach((expectedItem, index) => {
        if (typeof expectedItem === 'object' && expectedItem !== null) {
          const actualItem = actual[index];

          Object.entries(expectedItem).forEach(([key, value]) => {
            if (typeof value === 'object' && value !== null) {
              makeAssertion(
                actualItem[key],
                value,
                `${fieldPath}[${index}].${key}`,
              );
            } else {
              expect
                .soft(
                  actualItem[key],
                  `Field '${fieldPath}[${index}].${key}' value mismatch`,
                )
                .toEqual(value);
            }
          });
        } else {
          expect
            .soft(actual[index], `Field '${fieldPath}[${index}] value mismatch`)
            .toEqual(expectedItem);
        }
      });
    }
  } else if (expected === null) {
    expect.soft(actual, `Field '${fieldPath}' should be null`).toBeNull();
  } else if (typeof expected === 'object') {
    expect
      .soft(actual, `Field '${fieldPath}' should be an object`)
      .toBeDefined();

    if (actual) {
      Object.entries(expected).forEach(([key, value]) => {
        makeAssertion(
          (actual as Record<string, unknown>)[key],
          value,
          `${fieldPath}.${key}`,
        );
      });
    }
  } else {
    expect
      .soft(actual, `Field '${fieldPath}' value mismatch`)
      .toEqual(expected);
  }
}

function verifyNestedState(
  actual: unknown,
  expected: ExpectedValue,
  path: string[],
): void {
  if (expected === null) {
    makeAssertion(actual, null, createFieldPath(path));

    return;
  }

  if (Array.isArray(expected)) {
    makeAssertion(actual, expected, createFieldPath(path));

    return;
  }

  if (typeof expected === 'object' && 'assertType' in expected) {
    const assertion = expect.soft(
      actual,
      `Field '${createFieldPath(path)}' assertion failed`,
    );

    switch ((expected as SpecialAssertion).assertType) {
      case 'not.toBeNull':
        assertion.not.toBeNull();
        break;
      case 'toBeNull':
        assertion.toBeNull();
        break;
    }

    return;
  }

  if (typeof expected === 'object') {
    Object.entries(expected as Record<string, ExpectedValue>).forEach(
      ([key, value]) => {
        verifyNestedState(
          (actual as Record<string, unknown>)?.[key],
          value as ExpectedValue,
          [...path, key],
        );
      },
    );

    return;
  }

  makeAssertion(actual, expected, createFieldPath(path));
}

async function verifyOrganizationState(
  organizationId: string,
  expectedState: ExpectedState,
  customAssertions: Record<string, ExpectedValue>,
  maxRetries = 3,
): Promise<void> {
  let retries = 0;
  let assertionsPassed = false;

  await sleep(1000);

  while (retries < maxRetries && !assertionsPassed) {
    try {
      const { organization } = await service.getOrganization(organizationId);

      if (!organization) {
        throw new Error('Organization not found');
      }

      const modifiedExpectedState = JSON.parse(
        JSON.stringify(expectedState),
      ) as ExpectedState;

      Object.entries(customAssertions).forEach(([path, value]) => {
        const pathParts = path.split('.');

        if (pathParts.length === 1) {
          modifiedExpectedState[pathParts[0]] = value as ExpectedValue;
        } else {
          let current = modifiedExpectedState as Record<string, unknown>;

          for (let i = 0; i < pathParts.length - 1; i++) {
            const part = pathParts[i];
            const nextPart = pathParts[i + 1];

            if (!(part in current)) {
              current[part] = !isNaN(Number(nextPart)) ? [] : {};
            } else if (current[part] === null) {
              current[part] = !isNaN(Number(nextPart)) ? [] : {};
            }

            current = current[part] as Record<string, unknown>;
          }

          const lastPart = pathParts[pathParts.length - 1];

          current[lastPart] = value;
        }
      });

      Object.entries(modifiedExpectedState).forEach(([key, value]) => {
        verifyNestedState(
          organization[key as keyof Organization],
          value as ExpectedValue,
          [key],
        );
      });

      assertionsPassed = true;
    } catch (error) {
      retries++;

      if (retries < maxRetries) {
        await sleep(500);
      } else {
        throw error;
      }
    }
  }
}

const DEFAULT_ORGANIZATION_STATE: ExpectedState = {
  accountDetails: {
    churned: null,
    ltv: 0,
    onboarding: {
      status: 'NOT_APPLICABLE',
      comments: '',
      updatedAt: null,
    },
    renewalSummary: {
      arrForecast: null,
      maxArrForecast: null,
      renewalLikelihood: null,
      nextRenewalDate: null,
    },
  },
  contracts: null,
  description: '',
  domains: [],
  employees: 0,
  icon: '',
  industry: '',
  isCustomer: false,
  lastTouchpoint: {
    lastTouchPointAt: { assertType: 'not.toBeNull' },
    lastTouchPointTimelineEvent: { assertType: 'not.toBeNull' },
    lastTouchPointTimelineEventId: { assertType: 'not.toBeNull' },
    lastTouchPointType: { assertType: 'not.toBeNull' },
  },
  leadSource: '',
  locations: [],
  logo: '',
  owner: null,
  parentCompanies: [],
  public: false,
  relationship: '',
  socialMedia: [],
  subsidiaries: [],
  stage: '',
  valueProposition: '',
  yearFounded: null,
  website: '',
  tags: null,
};

describe('OrganizationsService - Integration Tests', () => {
  it('gets organizations', async () => {
    const { dashboardView_Organizations } = await service.getOrganizations({
      pagination: {
        page: 0,
        limit: 1000,
      },
    });

    expect(dashboardView_Organizations).toHaveProperty('content');
    expect(dashboardView_Organizations).toHaveProperty('totalElements');
    expect(dashboardView_Organizations).toHaveProperty('totalAvailable');

    const data = dashboardView_Organizations?.content;
    const totalElements = dashboardView_Organizations?.totalElements;
    const totalAvailable = dashboardView_Organizations?.totalAvailable;

    expect(data).toHaveLength(totalElements);
    expect(totalElements).toBeLessThanOrEqual(totalAvailable);
  });

  it('checks create empty organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await service.saveOrganization({
      input: { name: organization_name },
    });

    const customState = {
      ...DEFAULT_ORGANIZATION_STATE,
      name: organization_name,
    };

    await verifyOrganizationState(
      organization_Save.metadata.id,
      customState,
      {},
    );
  });

  it('adds tags to organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const organization_tag_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await service.saveOrganization({
      input: { name: organization_name },
    });

    await service.addTag({
      input: {
        organizationId: organization_Save.metadata.id,
        tag: { name: organization_tag_name },
      },
    });

    const customState = {
      ...DEFAULT_ORGANIZATION_STATE,
      name: organization_name,
      tags: [{ name: organization_tag_name }],
    };

    await verifyOrganizationState(
      organization_Save.metadata.id,
      customState,
      {},
    );
  });

  it('adds social to organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const organization_social_url = 'www.IT_' + crypto.randomUUID() + '.com';
    const { organization_Save } = await service.saveOrganization({
      input: { name: organization_name },
    });

    await service.addSocial({
      organizationId: organization_Save.metadata.id,
      input: {
        url: organization_social_url,
      },
    });

    const customState = {
      ...DEFAULT_ORGANIZATION_STATE,
      name: organization_name,
      socialMedia: [{ url: organization_social_url }],
    };

    await verifyOrganizationState(
      organization_Save.metadata.id,
      customState,
      {},
    );
  });

  it('adds subsidiary to organization', async () => {
    const parent_organization_name = 'IT_' + crypto.randomUUID();
    const subsidiary_organization_name = 'IT_' + crypto.randomUUID();

    const parent_organization = await service.saveOrganization({
      input: { name: parent_organization_name },
    });

    const subsidiary_organization = await service.saveOrganization({
      input: { name: subsidiary_organization_name },
    });

    await service.addSubsidiary({
      input: {
        organizationId: parent_organization.organization_Save.metadata.id,
        subsidiaryId: subsidiary_organization.organization_Save.metadata.id,
      },
    });

    const customState = {
      ...DEFAULT_ORGANIZATION_STATE,
      subsidiaries: [
        {
          organization: {
            name: subsidiary_organization_name,
          },
        },
      ],
      name: parent_organization_name,
    };

    await verifyOrganizationState(
      parent_organization.organization_Save.metadata.id,
      customState,
      {},
    );
  });
});
