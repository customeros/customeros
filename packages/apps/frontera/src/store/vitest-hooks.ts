import { afterAll } from 'vitest';

import { Transport } from './transport';
import { TagService } from './Tags/Tag.service';

const transport = new Transport();
// const organizationsService = OrganizationsService.getInstance(transport);
const tagService = TagService.getInstance(transport);

afterAll(async () => {
  // const where = {
  //   AND: [
  //     {
  //       filter: {
  //         property: 'UPDATED_AT',
  //         value: lastActiveAtUTC,
  //         operation: ComparisonOperator.Gte,
  //       },
  //     },
  //   ],
  // };
  //
  // const retrieved_organizations = await organizationsService.getOrganizations({
  //   pagination: { limit: 1000, page: 0 },
  //   sort: {
  //     by: 'LAST_TOUCHPOINT',
  //     caseSensitive: false,
  //     direction: SortingDirection.Desc,
  //   },
  // });

  // await tagService.updateTag({
  //   input: {
  //     id: 'b794b214-35c6-41c8-9c39-7814804ee427',
  //     name: 'b794b214-35c6-41c8-9c39-7814804ee427XXX',
  //   },
  // });
  await tagService
    .getTags()
    .then((res) =>
      res.tags.filter((tag) => tag.name.includes('IT_')).map((tag) => tag.id),
    );
});
