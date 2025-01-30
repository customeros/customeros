import { afterAll } from 'vitest';
import { OrganizationRepository } from '@infra/repositories/organization';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

import { TagService } from './Tags/__service__/Tag.service';

const organizationsRepository = OrganizationRepository.getInstance();
const tagService = TagService.getInstance();

afterAll(async () => {
  const tagIds = await tagService
    .getTags()
    .then((res) =>
      res.tags
        .filter((tag) => tag.name.includes('IT_'))
        .map((tag) => tag.metadata.id),
    );

  for (const tagId of tagIds) {
    await tagService.deleteTag({ id: tagId });
  }

  const organizationIds = Array.from(
    organizationsTestState.createdOrganizationIds,
  );

  for (const id of organizationIds) {
    try {
      await organizationsRepository.hideOrganizations({ ids: [id] });
    } catch (error) {
      console.error(`Failed to cleanup organization ${id}:`, error);
    }
  }
});
