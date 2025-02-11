import { afterAll } from 'vitest';
import { OrganizationRepository } from '@infra/repositories/organization';
import { ContactService } from '@store/Contacts/__service__/Contacts.service.ts';
import { contactsTestState } from '@store/Contacts/__tests__/contactsTestState.ts';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

import { TagService } from './Tags/__service__/Tag.service';

const organizationsRepository = OrganizationRepository.getInstance();
const contactService = ContactService.getInstance();
const tagService = TagService.getInstance();

afterAll(async () => {
  const tagIds = await tagService
    .getTags()
    .then((res) =>
      res.tags
        .filter((tag) => tag.name.includes('Vitest_'))
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
      console.error(`Failed to cleanup company ${id}:`, error);
    }
  }

  const contactIds = Array.from(contactsTestState.createdContactsIds);

  for (const id of contactIds) {
    try {
      await contactService.archiveContact({ contactId: id });
    } catch (error) {
      console.error(`Failed to cleanup contact ${id}:`, error);
    }
  }
});
