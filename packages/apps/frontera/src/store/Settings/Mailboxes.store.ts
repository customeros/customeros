import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import {
  override,
  computed,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import { MailstackBuyRequest } from '@shared/types/__generated__/graphql.types';

import { MailboxStore } from './Mailbox.store';
import { MailboxesService } from './__service__/Mailboxes/Mailboxes.service';

export class MailboxesStore extends SyncableGroup<
  MailstackBuyRequest,
  MailboxStore
> {
  private service: MailboxesService;
  domain: string = '';
  baseBundle: Set<string> = new Set();
  extendedBundle: Set<string> = new Set();
  invalidDomains: string[] = [];
  domainSuggestions: string[] = [];
  usernames: [string, string] = ['', ''];

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, MailboxStore);
    this.service = MailboxesService.getInstance(transport);

    makeObservable<MailboxesStore>(this, {
      channelName: override,
      domain: observable,
      baseBundle: observable,
      extendedBundle: observable,
      invalidDomains: observable,
      domainSuggestions: observable,
      usernames: observable,
      hasUsernames: computed,
      domainCount: computed,
      mailboxesCount: computed,
      usernamesCount: computed,
    });
  }

  get hasUsernames() {
    return this.usernames[0].length > 0 || this.usernames[1].length > 0;
  }

  get usernamesCount() {
    return this.usernames.filter((v) => v !== '').length;
  }

  get domainCount() {
    return this.baseBundle.size + this.extendedBundle.size;
  }

  get mailboxesCount() {
    return this.usernames.reduce(
      (acc, curr) => (curr.length ? this.domainCount : 0) + acc,
      0,
    );
  }

  public selectDomain(domain: string) {
    runInAction(() => {
      if (this.baseBundle.size < 5) {
        this.baseBundle.add(domain);
      } else {
        this.extendedBundle.add(domain);
      }

      this.domainSuggestions = this.domainSuggestions.filter(
        (d) => d !== domain,
      );
    });
  }

  public removeDomain(domain: string) {
    runInAction(() => {
      this.baseBundle.delete(domain);
      this.extendedBundle.delete(domain);
    });
  }

  async bootstrap() {
    try {
      const { mailstack_RegisteredBuyDomainsWithMailboxes: mailboxes } =
        await this.service.getRegisteredMailboxes();

      this.load(mailboxes as MailstackBuyRequest[], {
        getId: (data) => data.id,
      });

      runInAction(() => {
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  async getDomainSuggestions() {
    try {
      const { mailstack_DomainPurchaseSuggestions } =
        await this.service.getMailstackDomainsSuggestions({
          domain: this.domain,
        });

      runInAction(() => {
        this.domainSuggestions = mailstack_DomainPurchaseSuggestions;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  public setDomainName(domain: string) {
    runInAction(() => {
      this.domain = domain;
    });
  }

  public setUsername(index: 0 | 1, username: string) {
    runInAction(() => {
      this.usernames[index] = username;
    });
  }

  async getPaymentIntent() {
    try {
      const { mailstack_RegisterBuyDomainsWithMailboxes } =
        await this.service.createMailbox({
          domains: [...this.baseBundle, ...this.extendedBundle],
          amount: 100,
          usernames: this.usernames.filter((v) => v !== ''),
        });

      return mailstack_RegisterBuyDomainsWithMailboxes;
    } catch (err) {
      this.root.ui.toastError(
        'Failed processing the payment.',
        'get-payment-intent',
      );
    }
  }

  async validateDomains({
    onSuccess,
    onInvalid,
  }: { onSuccess?: () => void; onInvalid?: () => void } = {}) {
    try {
      runInAction(() => {
        this.isLoading = true;
      });

      const response = await this.service.getMailstackCheckUnavailableDomains({
        domains: [...this.baseBundle, ...this.extendedBundle],
      });

      runInAction(() => {
        this.invalidDomains = response.mailstack_CheckUnavailableDomains;

        if (this.invalidDomains.length === 0) {
          onSuccess?.();
        }

        if (this.invalidDomains.length > 0) {
          onInvalid?.();
        }
      });
    } catch (err) {
      this.root.ui.toastError('Failed validating domains', 'validate-domains');
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  get channelName() {
    return 'Mailboxes';
  }

  get persisterKey() {
    return 'Mailboxes';
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray<T extends MailboxesStore>(
    compute: (arr: MailboxStore[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }
}
