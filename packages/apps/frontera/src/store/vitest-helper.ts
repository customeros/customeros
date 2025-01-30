// import { expect } from 'vitest';
import { OrganizationRepository } from '@infra/repositories/organization';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

export class VitestHelper {
  static async createOrganizationForTest(
    organizationRepository: OrganizationRepository,
    input?: { input: { id?: string; name?: string; ownerId?: string } },
  ) {
    const organization_name = 'vitest-' + crypto.randomUUID();
    const organization_domain = 'vitest-' + crypto.randomUUID() + '.com';
    const { organization_Save } = await organizationRepository.saveOrganization(
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
