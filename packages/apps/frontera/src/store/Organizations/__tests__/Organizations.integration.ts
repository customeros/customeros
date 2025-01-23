import { it, expect, describe } from 'vitest';
import { VitestHelper } from '@store/vitest-helper.ts';

import {
  Currency,
  BilledType,
  OnboardingStatus,
  SortingDirection,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { Transport } from '../../../infra/transport';
import { UserService } from '../../Users/User.service';
import { ContractService } from '../../Contracts/Contract.service';
import { OrganizationsService } from '../__service__/Organizations.service';
import { ContractLineItemService } from '../../ContractLineItems/ContractLineItem.service';

const transport = new Transport();
const organizationsService = OrganizationsService.getInstance();
const contractService = ContractService.getInstance(transport);
const contractLineItemsService = ContractLineItemService.getInstance(transport);
const userService = UserService.getInstance(transport);

describe('OrganizationsService - Integration Tests', () => {
  it('gets organizations', async () => {
    const { dashboardView_Organizations } =
      await organizationsService.getOrganizations({
        pagination: {
          page: 0,
          limit: 10000,
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
    const { id, name } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const sleep = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));
    const maxRetries = 3;
    let retries = 0;
    let assertionsPassed = false;

    await sleep(1000);

    while (retries < maxRetries && !assertionsPassed) {
      try {
        const organization = await organizationsService.getOrganization(id);

        expect.soft(organization.churnedAt).toBeNull();
        expect.soft(organization.ltv).toBe(0);
        expect.soft(organization.onboardingStatus).toBe('NOT_APPLICABLE');
        expect.soft(organization.onboardingComments).toBe('');
        expect.soft(organization.onboardingStatusUpdatedAt).toBeNull();
        expect.soft(organization.renewalSummaryArrForecast).toBeNull();
        expect.soft(organization.renewalSummaryMaxArrForecast).toBeNull();
        expect.soft(organization.renewalSummaryRenewalLikelihood).toBeNull();
        expect.soft(organization?.renewalSummaryNextRenewalAt).toBeNull();
        expect.soft(organization?.contracts).toEqual([]);
        expect.soft(organization?.description).toBe('');
        expect.soft(organization?.employees).toEqual(0);
        expect.soft(organization?.iconUrl).toBe('');
        expect.soft(organization?.industryCode).toBeNull();
        expect(organization.lastTouchPointAt).not.toBeNull();
        expect.soft(organization.lastTouchPointType).not.toBeNull();
        expect.soft(organization?.leadSource).toBe('');
        expect.soft(organization?.locations).toEqual([]);
        expect.soft(organization?.logoUrl).toBe('');
        expect.soft(organization?.name).toBe(name);
        expect.soft(organization?.owner).toBeNull();
        expect.soft(organization?.parentId).toBeNull();
        expect.soft(organization?.parentName).toBeNull();
        expect.soft(organization?.public).toBe(false);
        expect.soft(organization?.relationship).toBe('PROSPECT');
        expect.soft(organization?.tags).toEqual([]);
        expect.soft(organization?.socialMedia).toEqual([]);
        expect.soft(organization?.subsidiaries).toEqual([]);

        expect.soft(organization?.stage).toBe('LEAD');
        expect.soft(organization?.yearFounded).toBeNull();
        expect.soft(organization?.website).toBe('');

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
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const organization_tag_name = 'Vitest_' + crypto.randomUUID();

    await organizationsService.addTag({
      input: {
        organizationId: id,
        tag: { name: organization_tag_name },
      },
    });

    let organization;

    organization = await organizationsService.getOrganization(id);
    expect(organization?.tags?.[0].name).toEqual(organization_tag_name);

    if (organization?.tags?.[0]?.metadata.id) {
      await organizationsService.removeTag({
        input: {
          organizationId: id,
          tag: { id: organization?.tags[0].metadata.id },
        },
      });

      organization = await organizationsService.getOrganization(id);
      expect(organization?.tags).toEqual([]);
    } else {
      throw new Error(
        'Tag removal failed: Organization or tag ID is undefined.',
      );
    }
  });

  it('adds social to organization', async () => {
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const organization_initial_social_url =
      'www.Vitest_' + crypto.randomUUID() + '.com';
    const { organization_AddSocial } = await organizationsService.addSocial({
      organizationId: id,
      input: {
        url: organization_initial_social_url,
      },
    });

    let organization;

    organization = await organizationsService.getOrganization(id);

    expect(organization?.socialMedia[0].url).toEqual(
      organization_initial_social_url,
    );

    const organization_subsequent_social_url =
      'www.Vitest_' + crypto.randomUUID() + '.com';

    await organizationsService.updateSocial({
      input: {
        id: organization_AddSocial.id,
        url: organization_subsequent_social_url,
      },
    });
    organization = await organizationsService.getOrganization(id);

    expect(organization?.socialMedia[0].url).toEqual(
      organization_subsequent_social_url,
    );

    await organizationsService.removeSocial({
      socialId: organization_AddSocial.id,
    });

    await new Promise((resolve) => setTimeout(resolve, 1000)); // waits for 1 second
    organization = await organizationsService.getOrganization(id);

    expect(organization?.socialMedia.length).toBe(0);
  });

  it('adds subsidiary to organization', async () => {
    const parent = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const subsidiary = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    await organizationsService.addSubsidiary({
      input: {
        organizationId: parent.id,
        subsidiaryId: subsidiary.id,
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
        organization = await organizationsService.getOrganization(parent.id);

        expect(organization?.subsidiaries[0]).toEqual(subsidiary.id);

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
      organizationId: parent.id,
      subsidiaryId: subsidiary.id,
    });

    while (retries < maxRetries && !assertionsPassed) {
      try {
        organization = await organizationsService.getOrganization(parent.id);

        expect(organization?.subsidiaries).toHaveLength(0);

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

    const { name, id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

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

    const hasName = (organizationName: string): boolean => {
      return (
        retrieved_organizations?.dashboardView_Organizations?.content?.some(
          (org) => org.name === organizationName,
        ) ?? false
      );
    };

    let organizationExistsInDashboard = hasName(name);

    expect(organizationExistsInDashboard).toBe(true);

    let archived_organizations =
      await organizationsService.getArchivedOrganizationsAfter({
        date: testStartDate,
      });

    expect(archived_organizations.organizations_HiddenAfter.includes(id)).toBe(
      false,
    );

    await organizationsService.hideOrganizations({
      ids: [id],
    });

    retrieved_organizations = await organizationsService.getOrganizations({
      pagination: { limit: 100, page: 0 },
      sort: {
        by: 'LAST_TOUCHPOINT',
        caseSensitive: false,
        direction: SortingDirection.Desc,
      },
    });

    organizationExistsInDashboard = hasName(name);

    expect(organizationExistsInDashboard).toBe(false);

    archived_organizations =
      await organizationsService.getArchivedOrganizationsAfter({
        date: testStartDate,
      });

    expect(archived_organizations.organizations_HiddenAfter.includes(id)).toBe(
      true,
    );
  });

  it.skip('updates onboarding status to organization', async () => {
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    let organization;

    organization = await organizationsService.getOrganization(id);
    expect(organization.onboardingStatus).toBe('NOT_APPLICABLE');
    await organizationsService.updateOnboardingStatus({
      input: {
        organizationId: id,
        status: OnboardingStatus.Stuck,
      },
    });

    organization = await organizationsService.getOrganization(id);
    expect(organization?.onboardingStatus).toBe('STUCK');
  });

  it('updates updateAllOpportunityRenewals', async () => {
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const contract_name = 'Vitest_' + crypto.randomUUID();
    const threeMonthsAgo = new Date();

    threeMonthsAgo.setMonth(threeMonthsAgo.getMonth() - 3);

    const { contract_Create } = await contractService.createContract({
      input: {
        organizationId: id,
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
        organizationId: id,
        renewalAdjustedRate: 53,
        renewalLikelihood: OpportunityRenewalLikelihood.HighRenewal,
      },
    });

    await organizationsService.updateAllOpportunityRenewals({
      input: {
        organizationId: id,
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
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    let organization = await organizationsService.getOrganization(id);

    expect(organization?.owner).toBeNull();

    const user = await userService.getUsers({
      pagination: { limit: 1000, page: 0 },
    });

    await organizationsService.saveOrganization({
      input: {
        id,
        ownerId: user.users.content[0].id,
      },
    });

    organization = await organizationsService.getOrganization(id);
    expect(organization?.owner?.id).toBe(user.users.content[0].id);
  });
});
