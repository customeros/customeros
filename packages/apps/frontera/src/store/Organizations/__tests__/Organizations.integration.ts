import { it, expect, describe } from 'vitest';

import { Transport } from '../../transport';
import { OrganizationsService } from '../__service__/Organizations.service';

const transport = new Transport();
const service = OrganizationsService.getInstance(transport);

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

  it('check create empty organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await service.saveOrganization({
      input: { name: organization_name },
    });

    const organization = await service.getOrganization(
      organization_Save.metadata.id,
    );

    expect.soft(organization.organization?.accountDetails?.churned).toBeNull;
    expect.soft(organization.organization?.accountDetails?.ltv).toBe(0);
    expect
      .soft(organization.organization?.accountDetails?.onboarding?.status)
      .toBe('NOT_APPLICABLE');
    expect
      .soft(organization.organization?.accountDetails?.onboarding?.comments)
      .toBe('');
    expect.soft(
      organization.organization?.accountDetails?.onboarding?.updatedAt,
    ).toBeNull;
    expect.soft(
      organization.organization?.accountDetails?.renewalSummary?.arrForecast,
    ).toBeNull;
    expect.soft(
      organization.organization?.accountDetails?.renewalSummary?.maxArrForecast,
    ).toBeNull;
    expect.soft(
      organization.organization?.accountDetails?.renewalSummary
        ?.renewalLikelihood,
    ).toBeNull;
    expect.soft(
      organization.organization?.accountDetails?.renewalSummary
        ?.nextRenewalDate,
    ).toBeNull;
    expect.soft(organization.organization?.contracts).toBeNull();
    expect.soft(organization.organization?.description).toBe('');
    expect.soft(organization.organization?.domains).toEqual([]);
    expect.soft(organization.organization?.employees).toEqual(0);
    expect.soft(organization.organization?.icon).toBe('');
    expect.soft(organization.organization?.industry).toBe('');
    expect.soft(organization.organization?.isCustomer).toBe(false);
    expect
      .soft(organization.organization?.lastTouchpoint?.lastTouchPointAt)
      .toBeNull();
    expect
      .soft(
        organization.organization?.lastTouchpoint?.lastTouchPointTimelineEvent,
      )
      .toBeNull();
    expect
      .soft(
        organization.organization?.lastTouchpoint
          ?.lastTouchPointTimelineEventId,
      )
      .toBeNull();
    expect
      .soft(organization.organization?.lastTouchpoint?.lastTouchPointType)
      .toBeNull();
    expect.soft(organization.organization?.leadSource).toBe('');
    expect.soft(organization.organization?.locations).toEqual([]);
    expect.soft(organization.organization?.logo).toBe('');
    expect.soft(organization.organization?.name).toBe(organization_name);
    expect.soft(organization.organization?.owner).toBeNull;
    expect.soft(organization.organization?.parentCompanies).toEqual([]);
    expect.soft(organization.organization?.public).toBe(false);
    expect.soft(organization.organization?.relationship).toBe('');
    expect.soft(organization.organization?.tags).toBeNull;
    expect.soft(organization.organization?.socialMedia).toEqual([]);
    expect.soft(organization.organization?.subsidiaries).toEqual([]);
    expect.soft(organization.organization?.stage).toBe('');
    expect.soft(organization.organization?.valueProposition).toBe('');
    expect.soft(organization.organization?.yearFounded).toBeNull();
    expect.soft(organization.organization?.website).toBe('');
  });
});
