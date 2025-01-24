import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import { GlobalCacheQuery } from './__service__/getGlobalCache.generated';
import { GlobalCacheService } from './__service__/GlobalCache.service.ts';
import { UserUpdateOnboardingDetailsMutationVariables } from './__service__/updateUserOnboardingDetails.generated';

export class GlobalCacheStore {
  value: GlobalCacheQuery['global_Cache'] | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;
  private service: GlobalCacheService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    this.service = GlobalCacheService.getInstance(transport);
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    await this.load();
  }

  async load() {
    try {
      this.isLoading = true;

      const response = await this.service.getGlobalCache();

      this.value = response.global_Cache;
      this.isBootstrapped = true;
    } catch (error) {
      this.error = (error as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  async updateOnboardingDetails(
    payload: Omit<
      Partial<UserUpdateOnboardingDetailsMutationVariables['input']>,
      'id'
    >,
    options?: { onSuccess?: () => void },
  ) {
    if (this.root.demoMode) return;

    if (this.value?.user?.onboarding) {
      const onboarding = this.value.user.onboarding;

      this.value.user.onboarding = {
        showOnboardingPage:
          payload?.showOnboardingPage ?? onboarding.showOnboardingPage,
        onboardingInboundStepCompleted:
          payload?.onboardingInboundStepCompleted ??
          onboarding.onboardingInboundStepCompleted,
        onboardingOutboundStepCompleted:
          payload?.onboardingOutboundStepCompleted ??
          onboarding.onboardingOutboundStepCompleted,
        onboardingCrmStepCompleted:
          payload?.onboardingCrmStepCompleted ??
          onboarding.onboardingCrmStepCompleted,
        onboardingMailstackStepCompleted:
          payload?.onboardingMailstackStepCompleted ??
          onboarding.onboardingMailstackStepCompleted,
      };
    }

    if (!this.value?.user || !this.value.user.id) {
      this.root.ui.toastError('User not found', 'user-not-found');

      return;
    }

    try {
      await this.service.updateOnboardingDetails({
        input: {
          ...this.value.user.onboarding,
          ...payload,
          id: this.value.user.id,
        },
      });

      runInAction(() => {
        options?.onSuccess?.();
      });
    } catch (error) {
      this.error = (error as Error)?.message;
    }
  }
}
