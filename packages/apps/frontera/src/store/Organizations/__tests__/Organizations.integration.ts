import { it, expect, describe } from 'vitest';

import { Transport } from '../../transport';
import { OrganizationsService } from '../__service__/Organizations.service';

const transport = new Transport();

transport.setHeaders({
  'X-OPENLINE-API-KEY': import.meta.env.VITE_TEST_API_KEY as string,
  'X-OPENLINE-USERNAME': import.meta.env.VITE_TEST_USERNAME as string,
});

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
});
