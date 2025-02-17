import { RootStore } from '@store/root';
import { TagDatum } from '@store/Tags/Tag.store';
import { TagService as TagRepo } from '@store/Tags/__service__/Tag.service';

import { unwrap } from '@utils/unwrap';

export class TagService {
  private tagRepo = TagRepo.getInstance();
  private store = RootStore.getInstance();

  constructor() {}

  public async createTag(
    payload: Partial<TagDatum>,
    options?: { onSuccess?: (serverId: string) => void },
  ) {
    if (!payload.name) {
      console.error('TagService.createTag: name is required');

      return;
    }

    const [res, err] = await unwrap(
      this.tagRepo.createTag({
        input: {
          name: payload.name,
          entityType: payload.entityType,
          colorCode: payload.colorCode,
        },
      }),
    );

    if (err) {
      console.error(err);

      return;
    }

    if (!res) {
      console.error('No response from createTag');

      return;
    }

    const serverId = res?.tag_Create.metadata.id;

    this.store.tags.createNew(res.tag_Create);
    this.store.tags.sync({
      action: 'APPEND',
      ids: [serverId],
    });

    if (serverId != null) {
      options?.onSuccess?.(serverId);
      this.store.ui.toastSuccess('Tag created', 'tag-created');
    }
  }
}
