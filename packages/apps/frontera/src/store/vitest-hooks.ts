import { afterAll } from 'vitest';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service.ts';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

import { Transport } from './transport';
import { TagService } from './Tags/Tag.service';

const transport = new Transport();
const organizationsService = OrganizationsService.getInstance(transport);
const tagService = TagService.getInstance(transport);

afterAll(async () => {
  const tagIds = await tagService
    .getTags()
    .then((res) =>
      res.tags.filter((tag) => tag.name.includes('IT_')).map((tag) => tag.id),
    );

  for (const tagId of tagIds) {
    await tagService.deleteTag({ id: tagId });
  }

  const organizationIds = Array.from(
    organizationsTestState.createdOrganizationIds,
  );

  for (const id of organizationIds) {
    try {
      await organizationsService.hideOrganizations({ ids: [id] });
    } catch (error) {
      console.error(`Failed to cleanup organization ${id}:`, error);
    }
  }
});
