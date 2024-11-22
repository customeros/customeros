import { it, expect, describe } from 'vitest';
import { VitestHelper } from '@store/vitest-helper.ts';

import {
  Currency,
  BilledType,
  OnboardingStatus,
  SortingDirection,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { Transport } from '../../transport';
import { UserService } from '../../Users/User.service';
import { ContractService } from '../../Contracts/Contract.service';
import { OrganizationsService } from '../__service__/Organizations.service';
import { ContractLineItemService } from '../../ContractLineItems/ContractLineItem.service';

const transport = new Transport();
const organizationsService = OrganizationsService.getInstance(transport);
const contractService = ContractService.getInstance(transport);
const contractLineItemsService = ContractLineItemService.getInstance(transport);
const userService = UserService.getInstance(transport);

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
    const { organization_Save, organization_name } =
      await VitestHelper.createOrganizationForTest(organizationsService);

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
        expect(
          organization.organization?.lastTouchpoint?.lastTouchPointAt,
        ).not.toBeNull();
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
    const { organization_Save } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const organization_tag_name = 'IT_' + crypto.randomUUID();

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

    if (organization?.organization?.tags?.[0]?.metadata.id) {
      await organizationsService.removeTag({
        input: {
          organizationId: organization_Save.metadata.id,
          tag: { id: organization.organization.tags[0].metadata.id },
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
    const { organization_Save } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const organization_initial_social_url =
      'www.IT_' + crypto.randomUUID() + '.com';
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
        id: organization_AddSocial.id,
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
    const parent_organization = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const {
      organization_Save: subsidiary_organization,
      organization_name: subsidiary_organization_name,
    } = await VitestHelper.createOrganizationForTest(organizationsService);

    await organizationsService.addSubsidiary({
      input: {
        organizationId: parent_organization.organization_Save.metadata.id,
        subsidiaryId: subsidiary_organization.metadata.id,
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
      subsidiaryId: subsidiary_organization.metadata.id,
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

    const {
      organization_name: organization_name,
      organization_Save: new_organization,
    } = await VitestHelper.createOrganizationForTest(organizationsService);

    const sleep = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));

    await sleep(500);

    let retrieved_organizations = await organizationsService.getOrganizations({
      pagination: { limit: 100, page: 0 },
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
        new_organization.metadata.id,
      ),
    ).toBe(false);

    await organizationsService.hideOrganizations({
      ids: [new_organization.metadata.id],
    });

    retrieved_organizations = await organizationsService.getOrganizations({
      pagination: { limit: 100, page: 0 },
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
        new_organization.metadata.id,
      ),
    ).toBe(true);
  });

  it('updates onboarding status to organization', async () => {
    const { organization_Save } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

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

  it('updates updateAllOpportunityRenewals', async () => {
    const { organization_Save } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const contract_name = 'IT_' + crypto.randomUUID();
    const threeMonthsAgo = new Date();

    threeMonthsAgo.setMonth(threeMonthsAgo.getMonth() - 3);

    const { contract_Create } = await contractService.createContract({
      input: {
        organizationId: organization_Save.metadata.id,
        committedPeriodInMonths: 3,
        currency: Currency.Usd,
        name: contract_name,
        serviceStarted: threeMonthsAgo,
      },
    });

    await contractLineItemsService.createContractLineItem({
      input: {
        billingCycle: BilledType.Monthly,
        contractId: contract_Create.metadata.id,
        description: 'SLI_001',
        price: 4,
        quantity: 3,
        serviceEnded: null,
        serviceStarted: threeMonthsAgo,
        tax: { taxRate: 5 },
      },
    });

    await organizationsService.updateAllOpportunityRenewals({
      input: {
        organizationId: organization_Save.metadata.id,
        renewalAdjustedRate: 53,
        renewalLikelihood: OpportunityRenewalLikelihood.HighRenewal,
      },
    });

    await organizationsService.updateAllOpportunityRenewals({
      input: {
        organizationId: organization_Save.metadata.id,
        renewalAdjustedRate: 53,
        renewalLikelihood: OpportunityRenewalLikelihood.HighRenewal,
      },
    });

    const updatedContract = await contractService
      .getContracts({
        pagination: { limit: 1000, page: 0 },
      })
      .then(
        (res) =>
          res.contracts.content.filter(
            (c) => c.metadata.id === contract_Create.metadata.id,
          )[0],
      );

    expect
      .soft(updatedContract.opportunities?.[0].renewalLikelihood)
      .toBe(OpportunityRenewalLikelihood.HighRenewal);
    expect
      .soft(updatedContract.opportunities?.[0].renewalAdjustedRate)
      .toBe(53);
    expect.soft(updatedContract.opportunities?.[0].amount).toBe(76.32);
    // );
  });

  it('adds updates owner of the organization', async () => {
    const { organization_Save } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    let organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );

    expect(organization.organization?.owner).toBeNull();

    const user = await userService.getUsers({
      pagination: { limit: 1000, page: 0 },
    });

    await organizationsService.saveOrganization({
      input: {
        id: organization_Save.metadata.id,
        ownerId: user.users.content[0].id,
      },
    });

    organization = await organizationsService.getOrganization(
      organization_Save.metadata.id,
    );
    expect(organization.organization?.owner?.id).toBe(user.users.content[0].id);
  });
});
