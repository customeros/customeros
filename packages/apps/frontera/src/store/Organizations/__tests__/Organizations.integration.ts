import { it, expect, describe } from 'vitest';

import { OnboardingStatus, SortingDirection } from '@graphql/types';

import { Transport } from '../../transport';
import { trackOrganization } from './organizationsTestState';
import { OrganizationsService } from '../__service__/Organizations.service';

const transport = new Transport();
const organizationsService = OrganizationsService.getInstance(transport);

describe('OrganizationsService - Integration Tests', () => {
  it('gets organizations', async () => {
    const { dashboardView_Organizations } =
      await organizationsService.getOrganizations({
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
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    trackOrganization(organization_Save.metadata.id);

    const sleep = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));
    const maxRetries = 3;
    let retries = 0;
    let organization;
    let assertionsPassed = false;

    await sleep(1000);

    while (retries < maxRetries && !assertionsPassed) {
      try {
        organization = await organizationsService.getOrganization(
          organization_Save.metadata.id,
        );

        expect
          .soft(organization.organization?.accountDetails?.churned)
          .toBeNull();
        expect.soft(organization.organization?.accountDetails?.ltv).toBe(0);
        expect
          .soft(organization.organization?.accountDetails?.onboarding?.status)
          .toBe('NOT_APPLICABLE');
        expect
          .soft(organization.organization?.accountDetails?.onboarding?.comments)
          .toBe('');
        expect
          .soft(
            organization.organization?.accountDetails?.onboarding?.updatedAt,
          )
          .toBeNull();
        expect
          .soft(
            organization.organization?.accountDetails?.renewalSummary
              ?.arrForecast,
          )
          .toBeNull();
        expect
          .soft(
            organization.organization?.accountDetails?.renewalSummary
              ?.maxArrForecast,
          )
          .toBeNull();
        expect
          .soft(
            organization.organization?.accountDetails?.renewalSummary
              ?.renewalLikelihood,
          )
          .toBeNull();
        expect
          .soft(
            organization.organization?.accountDetails?.renewalSummary
              ?.nextRenewalDate,
          )
          .toBeNull();
        expect.soft(organization.organization?.contracts).toBeNull();
        expect.soft(organization.organization?.description).toBe('');
        expect.soft(organization.organization?.domains).toEqual([]);
        expect.soft(organization.organization?.employees).toEqual(0);
        expect.soft(organization.organization?.icon).toBe('');
        expect.soft(organization.organization?.industry).toBe('');
        expect.soft(organization.organization?.isCustomer).toBe(false);
        expect
          .soft(organization.organization?.lastTouchpoint?.lastTouchPointAt)
          .not.toBeNull();
        expect
          .soft(
            organization.organization?.lastTouchpoint
              ?.lastTouchPointTimelineEvent,
          )
          .not.toBeNull();
        expect
          .soft(
            organization.organization?.lastTouchpoint
              ?.lastTouchPointTimelineEventId,
          )
          .not.toBeNull();
        expect
          .soft(organization.organization?.lastTouchpoint?.lastTouchPointType)
          .not.toBeNull();
        expect.soft(organization.organization?.leadSource).toBe('');
        expect.soft(organization.organization?.locations).toEqual([]);
        expect.soft(organization.organization?.logo).toBe('');
        expect.soft(organization.organization?.name).toBe(organization_name);
        expect.soft(organization.organization?.owner).toBeNull();
        expect.soft(organization.organization?.parentCompanies).toEqual([]);
        expect.soft(organization.organization?.public).toBe(false);
        expect.soft(organization.organization?.relationship).toBe('');
        expect.soft(organization.organization?.tags).toBeNull();
        expect.soft(organization.organization?.socialMedia).toEqual([]);
        expect.soft(organization.organization?.subsidiaries).toEqual([]);
        expect.soft(organization.organization?.stage).toBe('');
        expect.soft(organization.organization?.valueProposition).toBe('');
        expect.soft(organization.organization?.yearFounded).toBeNull();
        expect.soft(organization.organization?.website).toBe('');
        expect.soft(organization.organization?.website).toBe('');

        assertionsPassed = true;
      } catch (error) {
        retries++;

        if (retries < maxRetries) {
          await sleep(1000);
        } else {
          throw error;
        }
      }
    }
  });

  it('adds tags to organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const organization_tag_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    trackOrganization(organization_Save.metadata.id);

    await organizationsService.addTag({
      input: {
        organizationId: organization_Save.metadata.id,
        tag: { name: organization_tag_name },
      },
    });

    let organization;

    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );
    expect(organization.organization?.tags?.[0].name).toEqual(
      organization_tag_name,
    );

    if (organization?.organization?.tags?.[0]?.id) {
      await organizationsService.removeTag({
        input: {
          organizationId: organization_Save.metadata.id,
          tag: { id: organization.organization.tags[0].id },
        },
      });

      organization = await organizationsService.getOrganization(
        organization_Save.metadata.id,
      );
      expect(organization.organization?.tags).toBeNull();
    } else {
      throw new Error(
        'Tag removal failed: Organization or tag ID is undefined.',
      );
    }
  });

  it('adds social to organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const organization_initial_social_url =
      'www.IT_' + crypto.randomUUID() + '.com';

    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    trackOrganization(organization_Save.metadata.id);

    const { organization_AddSocial } = await organizationsService.addSocial({
      organizationId: organization_Save.metadata.id,
      input: {
        url: organization_initial_social_url,
      },
    });

    let organization;

    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );

    expect(organization.organization?.socialMedia[0].url).toEqual(
      organization_initial_social_url,
    );

    const organization_subsequent_social_url =
      'www.IT_' + crypto.randomUUID() + '.com';

    await organizationsService.updateSocial({
      input: {
        id: organization_Save.metadata.id,
        url: organization_subsequent_social_url,
      },
    });
    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );

    expect(organization.organization?.socialMedia[0].url).toEqual(
      organization_subsequent_social_url,
    );

    await organizationsService.removeSocial({
      socialId: organization_AddSocial.id,
    });

    await new Promise((resolve) => setTimeout(resolve, 1000)); // waits for 1 second
    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );

    expect(organization.organization?.socialMedia.length).toBe(0);
  });

  it('adds subsidiary to organization', async () => {
    const parent_organization_name = 'IT_' + crypto.randomUUID();
    const subsidiary_organization_name = 'IT_' + crypto.randomUUID();
    const parent_organization = await organizationsService.saveOrganization({
      input: { name: parent_organization_name },
    });

    trackOrganization(parent_organization.organization_Save.metadata.id);

    const subsidiary_organization = await organizationsService.saveOrganization(
      {
        input: { name: subsidiary_organization_name },
      },
    );

    trackOrganization(subsidiary_organization.organization_Save.metadata.id);

    await organizationsService.addSubsidiary({
      input: {
        organizationId: parent_organization.organization_Save.metadata.id,
        subsidiaryId: subsidiary_organization.organization_Save.metadata.id,
      },
    });

    const sleep = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));
    const maxRetries = 3;
    let retries = 0;
    let organization;
    let assertionsPassed = false;

    await sleep(500);

    while (retries < maxRetries && !assertionsPassed) {
      try {
        organization = await organizationsService.getOrganization(
          parent_organization.organization_Save.metadata.id,
        );

        expect(
          organization.organization?.subsidiaries[0].organization.name,
        ).toEqual(subsidiary_organization_name);

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

    await organizationsService.removeSubsidiary({
      organizationId: parent_organization.organization_Save.metadata.id,
      subsidiaryId: subsidiary_organization.organization_Save.metadata.id,
    });

    while (retries < maxRetries && !assertionsPassed) {
      try {
        organization = await organizationsService.getOrganization(
          parent_organization.organization_Save.metadata.id,
        );

        expect(
          organization.organization?.subsidiaries[0].organization.name,
        ).toBeNull();

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
  });

  it('retrieve archived organizations', async () => {
    const testStartDate = new Date().toISOString();

    const organization_name = 'IT_' + crypto.randomUUID();
    const new_organization = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    trackOrganization(new_organization.organization_Save.metadata.id);

    const sleep = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));

    await sleep(500);

    let retrieved_organizations = await organizationsService.getOrganizations({
      pagination: { limit: 1000, page: 0 },
      sort: {
        by: 'LAST_TOUCHPOINT',
        caseSensitive: false,
        direction: SortingDirection.Desc,
      },
    });

    let hasName = (organizationName: string): boolean => {
      return (
        retrieved_organizations?.dashboardView_Organizations?.content?.some(
          (org) => org.name === organizationName,
        ) ?? false
      );
    };
    let organizationExistsInDashboard = hasName(organization_name);

    expect(organizationExistsInDashboard).toBe(true);

    let archived_organizations =
      await organizationsService.getArchivedOrganizationsAfter({
        date: testStartDate,
      });

    expect(
      archived_organizations.organizations_HiddenAfter.includes(
        new_organization.organization_Save.metadata.id,
      ),
    ).toBe(false);

    await organizationsService.hideOrganizations({
      ids: [new_organization.organization_Save.metadata.id],
    });

    retrieved_organizations = await organizationsService.getOrganizations({
      pagination: { limit: 1000, page: 0 },
      sort: {
        by: 'LAST_TOUCHPOINT',
        caseSensitive: false,
        direction: SortingDirection.Desc,
      },
    });

    hasName = (organizationName: string): boolean => {
      return (
        retrieved_organizations?.dashboardView_Organizations?.content?.some(
          (org) => org.name === organizationName,
        ) ?? false
      );
    };
    organizationExistsInDashboard = hasName(organization_name);

    expect(organizationExistsInDashboard).toBe(false);

    archived_organizations =
      await organizationsService.getArchivedOrganizationsAfter({
        date: testStartDate,
      });

    expect(
      archived_organizations.organizations_HiddenAfter.includes(
        new_organization.organization_Save.metadata.id,
      ),
    ).toBe(true);
  });

  it('updates onboarding status to organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();

    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    trackOrganization(organization_Save.metadata.id);

    let organization;

    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );
    expect(organization.organization?.accountDetails?.onboarding?.status).toBe(
      'NOT_APPLICABLE',
    );
    await organizationsService.updateOnboardingStatus({
      input: {
        organizationId: organization_Save.metadata.id,
        status: OnboardingStatus.Stuck,
      },
    });

    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );
    expect(organization.organization?.accountDetails?.onboarding?.status).toBe(
      'STUCK',
    );
  });
});
