import { Tracer } from '@infra/tracer';
import { CommonRepository } from '@infra/repositories/common';

import { unwrap } from '@utils/unwrap';

export class CommonService {
  private repo = CommonRepository.getInstance();

  constructor() {}

  async getSlackChannels() {
    const span = Tracer.span('CommonService.getSlackChannels');
    const [data, err] = await unwrap(this.repo.getSlackChannels());

    if (err) {
      console.error(
        'CommonService.getSlackChannels: Error fetching slack channels',
        err,
      );
    }

    span.end();

    return data ? data.slackChannelsWithBot : [];
  }
}
