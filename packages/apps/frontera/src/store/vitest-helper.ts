// import { expect } from 'vitest';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service.ts';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

export class VitestHelper {
  static async createOrganizationForTest(
    organizationsService: OrganizationsService,
    input?: { input: { id?: string; name?: string; ownerId?: string } },
  ) {
    const organization_name = 'vitest-' + crypto.randomUUID();
    const organization_domain = 'vitest-' + crypto.randomUUID() + '.com';
    const { organization_Save } = await organizationsService.saveOrganization(
      input || {
        input: { name: organization_name, domains: [organization_domain] },
      },
    );
    const { metadata } = organization_Save;

    trackOrganization(metadata.id);
    // console.info(
    //   `\nOrganization ${organization_name} was created for test "${
    //     expect.getState().currentTestName
    //   }": `,
    // );

    return { id: metadata.id, name: organization_name };
  }
}

export const trackOrganization = (organizationId: string) => {
  organizationsTestState.createdOrganizationIds.add(organizationId);
};
