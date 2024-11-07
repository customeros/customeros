import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { UserStore } from '@store/Users/User.store';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { makeAutoSyncableGroup } from '@store/group-store';

import {
  Currency,
  DataSource,
  Opportunity,
  InternalType,
  InternalStage,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { OpportunitiesService } from './__services__/Opportunities.service';

export class OpportunityStore implements Store<Opportunity> {
  value: Opportunity = makeDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Opportunity>();
  update = makeAutoSyncable.update<Opportunity>();
  private service: OpportunitiesService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = OpportunitiesService.getInstance(transport);

    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Opportunity',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
  }

  get id() {
    return this.value.metadata.id;
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }

  get organization() {
    const organizationId = this.value.organization?.metadata.id;

    if (!organizationId) return null;

    return this.root.organizations.value.get(organizationId);
  }

  get externalStage() {
    const stages = this.root.settings.tenant.value?.opportunityStages;

    const stage = stages?.find((s) => s.value === this.value.externalStage);

    if (!stage) return null;

    return stage;
  }

  get owner() {
    const ownerId = this.value.owner?.id;

    if (!ownerId) return null;

    return this.root.users.value.get(ownerId);
  }

  async invalidate() {
    try {
      this.isLoading = true;

      const { opportunity } = await this.service.getOpportunity({
        id: this.id,
      });

      this.load(opportunity as Opportunity);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateProperty(property: keyof Opportunity) {
    try {
      this.isLoading = true;
      await this.service.saveOpportunity({
        input: {
          opportunityId: this.id,
          [property]: this.value[property],
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateOpportunityRenewal() {
    try {
      this.isLoading = true;

      const input = {
        opportunityId: this.id,
        comments: this.value.comments || '',
        renewalAdjustedRate: this.value.renewalAdjustedRate,
        renewalLikelihood: this.value.renewalLikelihood,
      };

      await this.service.updateOpportunityRenewal({
        input,
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        setTimeout(() => {
          this.invalidate();
        }, 1000);
      });
    }
  }

  private async updateOpportunityExternalStage(externalStage: string) {
    try {
      this.isLoading = true;
      await this.service.saveOpportunity({
        input: {
          opportunityId: this.id,
          externalStage,
          internalStage: InternalStage.Open,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateOpportunityCloseLost() {
    try {
      this.isLoading = true;
      await this.service.saveOpportunity({
        input: {
          opportunityId: this.id,
          externalStage: '',
          internalStage: InternalStage.ClosedLost,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateOpportunityCloseWon() {
    try {
      this.isLoading = true;
      await this.service.saveOpportunity({
        input: {
          opportunityId: this.id,
          externalStage: '',
          internalStage: InternalStage.ClosedWon,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateOpportunityOwner(userId: string) {
    try {
      this.isLoading = true;
      await this.service.saveOpportunity({
        input: {
          opportunityId: this.id,
          ownerId: userId,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;
    const value = diff?.val;

    match(path)
      .with(['externalStage'], () => {
        this.updateOpportunityExternalStage(value as string);
      })
      .with(['internalStage'], () => {
        match(value)
          .with(InternalStage.ClosedLost, () => {
            this.updateOpportunityCloseLost();
          })
          .with(InternalStage.ClosedWon, () => {
            this.updateOpportunityCloseWon();
          });
      })
      .with(['renewalLikelihood'], () => {
        this.updateOpportunityRenewal();
      })
      .with(['renewalAdjustedRate'], () => {
        this.updateOpportunityRenewal();
      })
      .with(['owner', ...P.array()], () => {
        if (typeof value === 'string') {
          this.updateOpportunityOwner(value as string);

          return;
        }

        if (value?.id) {
          this.updateOpportunityOwner(value.id);
        }
      })
      .otherwise(() => {
        const property = path?.[0] as keyof Opportunity;

        property && this.updateProperty(property);
      });
  }

  async saveProperty(property: keyof Opportunity) {
    this.updateProperty(property);
  }
}

const makeDefaultValue = (): Opportunity => ({
  metadata: {
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: '',
    source: DataSource.Na,
    sourceOfTruth: DataSource.Na,
    appSource: '',
  },
  amount: 0,
  appSource: '',
  comments: '',
  currency: Currency.Usd,
  createdAt: '',
  externalLinks: [],
  externalStage: '',
  externalType: '',
  generalNotes: '',
  id: '',
  likelihoodRate: 0,
  internalStage: InternalStage.Open,
  internalType: InternalType.Nbo,
  maxAmount: 0,
  name: '',
  nextSteps: '',
  owner: UserStore.getDefaultValue(),
  stageLastUpdated: '',
  renewalAdjustedRate: 0,
  renewalApproved: false,
  renewalLikelihood: OpportunityRenewalLikelihood.ZeroRenewal,
  renewalUpdatedByUserId: '',
  source: DataSource.Na,
  sourceOfTruth: DataSource.Na,
  updatedAt: new Date().toISOString(),
});
