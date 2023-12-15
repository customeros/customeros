import { SelectOption } from '@ui/utils/types';
import {
  OnboardingStatus,
  OnboardingDetails,
  OnboardingStatusInput,
} from '@graphql/types';

import { options } from './util';

export interface OnboardingStatusForm {
  comments: string;
  status: SelectOption<OnboardingStatus>;
}

export class OnboardingStatusDto {
  public comments: string;
  public status: SelectOption<OnboardingStatus>;

  constructor(data?: OnboardingDetails | null) {
    this.comments = data?.comments ?? '';
    this.status =
      options.find((option) => option.value === data?.status) ?? options[0];
  }

  static toForm(data?: OnboardingDetails | null): OnboardingStatusForm {
    return new OnboardingStatusDto(data);
  }

  static toPayload(
    input: OnboardingStatusForm & { id: string },
  ): OnboardingStatusInput {
    return {
      organizationId: input.id,
      comments: input.comments,
      status: input.status.value,
    };
  }
}
