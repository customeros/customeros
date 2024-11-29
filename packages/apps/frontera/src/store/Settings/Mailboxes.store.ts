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

import { validateUrl } from '@utils/url';

import { MailboxStore, type Mailbox } from './Mailbox.store';
import { MailboxesService } from './__service__/Mailboxes/Mailboxes.service';

export class MailboxesStore extends SyncableGroup<Mailbox, MailboxStore> {
  private service: MailboxesService;
  domain: string = '';
  baseBundle: Set<string> = new Set();
  extendedBundle: Set<string> = new Set();
  invalidDomains: string[] = [];
  domainSuggestions: string[] = [];
  redirectUrl: string = '';
  usernames: [string, string] = ['', ''];
  invalidUsernames: [string, string] = ['', ''];
  invalidRedirectUrl: string = '';
  invalidBaseBundle: string = '';

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
      redirectUrl: observable,
      hasUsernames: computed,
      domainCount: computed,
      mailboxesCount: computed,
      usernamesCount: computed,
      invalidRedirectUrl: observable,
      invalidUsernames: observable,
      invalidBaseBundle: observable,
      totalAmount: computed,
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

  get totalAmount() {
    // multiply by 100 to convert to cents (required by stripe)
    return parseFloat(
      ((199.99 + this.extendedBundle.size * 18.99) * 100).toFixed(2),
    );
  }

  public resetBuyFlow() {
    this.baseBundle.clear();
    this.extendedBundle.clear();
    this.usernames = ['', ''];
    this.domain = '';
    this.invalidDomains = [];
    this.domainSuggestions = [];
    this.redirectUrl = '';
    this.invalidUsernames = ['', ''];
    this.invalidBaseBundle = '';
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

      this.validateBaseBundle();
    });
  }

  public removeDomain(domain: string) {
    runInAction(() => {
      const newSet = new Set([...this.baseBundle, ...this.extendedBundle]);

      newSet.delete(domain);

      const newArr = Array.from(newSet);

      this.baseBundle = new Set(newArr.splice(0, 5));
      this.extendedBundle = new Set(newArr);

      this.validateBaseBundle();
      this.invalidDomains = this.invalidDomains.filter((d) => d !== domain);
    });
  }

  public validateRedirectUrl = () => {
    const isValidUrl = validateUrl(this.redirectUrl);
    let valid = false;

    runInAction(() => {
      if (this.redirectUrl.length === 0) {
        this.invalidRedirectUrl = 'Your domains need a destination';

        valid = false;

        return;
      }

      if (!isValidUrl) {
        this.invalidRedirectUrl = 'Invalid URL';

        valid = false;

        return;
      }

      valid = true;
      this.invalidRedirectUrl = '';
    });

    return valid;
  };

  public validateBaseBundle = () => {
    let valid = false;

    runInAction(() => {
      const count = this.baseBundle.size;

      if (count < 5) {
        this.invalidBaseBundle = `Please add ${5 - count} more ${
          count > 1 ? 'domains' : 'domain'
        }`;

        valid = false;

        return;
      }

      this.invalidBaseBundle = '';
      valid = true;
    });

    return valid;
  };

  public validateUsernames = () => {
    let valid = false;

    runInAction(() => {
      const [a, b] = this.usernames;

      if (a.length === 0 || b.length === 0) {
        this.invalidUsernames = [
          !a.length ? 'Houston we have a blank...' : '',
          !b.length ? 'Houston we have a blank...' : '',
        ];

        valid = false;

        return;
      }

      if (a.length > 0 && b.length > 0 && a === b) {
        this.invalidUsernames[0] = 'This username is already used';
        this.invalidUsernames[1] = 'This username is already used';

        valid = false;

        return;
      }

      this.invalidUsernames = ['', ''];
      valid = true;
    });

    return valid;
  };

  public async validateBuy({ onSuccess }: { onSuccess?: () => void }) {
    // sync validations
    const valid = [
      this.validateRedirectUrl(),
      this.validateUsernames(),
      this.validateBaseBundle(),
    ].every(Boolean);

    if (!valid) return;

    // async validations
    await this.validateDomains({ onSuccess });
  }

  async bootstrap() {
    try {
      const { mailstack_Mailboxes: mailboxes } =
        await this.service.getMailboxes();

      this.load(mailboxes, {
        getId: (data) => data.mailbox,
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
      this.isLoading = true;

      const { mailstack_DomainPurchaseSuggestions } =
        await this.service.getDomainSuggestions({
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

  public setRedirectUrl(url: string) {
    runInAction(() => {
      this.redirectUrl = url;
    });
  }

  public setUsername(index: 0 | 1, username: string) {
    runInAction(() => {
      this.usernames[index] = username;
    });
  }

  async getPaymentIntent() {
    try {
      const { mailstack_GetPaymentIntent } =
        await this.service.getPaymentIntent({
          domains: [...this.baseBundle, ...this.extendedBundle],
          amount: this.totalAmount,
          usernames: this.usernames.filter((v) => v !== ''),
        });

      return mailstack_GetPaymentIntent;
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

      const response = await this.service.validateDomains({
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

  async buyDomains(paymentIntentId: string) {
    try {
      await this.service.buyDomains({
        test: true,
        paymentIntentId,
        domains: [...this.baseBundle, ...this.extendedBundle],
        amount: this.totalAmount,
        username: this.usernames.filter((v) => v !== ''),
        redirectWebsite: this.redirectUrl,
      });
      this.bootstrap();
    } catch (err) {
      this.root.ui.toastError('Failed buying domains', 'buy-domains');
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
