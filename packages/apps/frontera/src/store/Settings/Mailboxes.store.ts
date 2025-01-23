import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
import { SyncableGroup } from '@store/syncable-group';
import {
  computed,
  override,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import { validateEmailLocalPart } from '@utils/email';
import { validateUrl, validateDomain as validateDomainStr } from '@utils/url';

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
  invalidDomain: string = '';
  dirty = new Map([
    ['username1', false],
    ['username2', false],
    ['redirectUrl', false],
  ]);

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
      invalidDomain: observable,
      dirty: observable,
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
    const baseBundlePrice = this.value.size > 0 ? 0 : 199.99;

    // multiply by 100 to convert to cents (required by stripe)
    return parseFloat(
      ((baseBundlePrice + this.extendedBundle.size * 18.99) * 100).toFixed(2),
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
    this.invalidRedirectUrl = '';
    this.invalidUsernames = ['', ''];
    this.invalidBaseBundle = '';
    this.invalidDomain = '';
    this.dirty.clear();
  }

  public setDirty(field: string) {
    runInAction(() => {
      this.dirty.set(field, true);
    });
  }

  public selectDomain(domain: string) {
    runInAction(() => {
      if (this.value.size === 0 && this.baseBundle.size < 5) {
        this.baseBundle.add(domain);
      } else {
        this.extendedBundle.add(domain);
      }

      this.domainSuggestions = this.domainSuggestions.filter(
        (d) => d !== domain,
      );

      this.invalidBaseBundle = '';
    });
  }

  public removeDomain(domain: string) {
    runInAction(() => {
      const newSet = new Set([...this.baseBundle, ...this.extendedBundle]);

      newSet.delete(domain);
      this.domainSuggestions.push(domain);

      const newArr = Array.from(newSet);

      if (this.value.size === 0) {
        this.baseBundle = new Set(newArr.splice(0, 5));
      }
      this.extendedBundle = new Set(newArr);

      this.invalidBaseBundle = '';
      this.invalidDomains = this.invalidDomains.filter((d) => d !== domain);
    });
  }

  public validateDomain = () => {
    let valid = false;

    runInAction(() => {
      if (this.domain === '') {
        this.invalidDomain = '';

        valid = true;

        return;
      }

      if (!validateDomainStr(this.domain)) {
        this.invalidDomain = 'Invalid domain format';

        valid = false;

        return;
      }

      this.invalidDomain = '';
      valid = true;
    });

    return valid;
  };

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
        this.invalidRedirectUrl = 'This URL appears to be invalid';

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
          5 - count > 1 ? 'domains' : 'domain'
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
    runInAction(() => {
      this.invalidUsernames = this.usernames.map((v, i, arr) => {
        if (v.length === 0) return 'Houston we have a blank...';
        if (!validateEmailLocalPart(v)) return 'Invalid username format';
        if (arr[i === 0 ? 1 : 0] === v) return 'This username is already used';

        return '';
      }) as [string, string];
    });

    return this.invalidUsernames.every((v) => v === '');
  };

  public async validateBuy({ onSuccess }: { onSuccess?: () => void }) {
    // sync validations
    const validators = [this.validateRedirectUrl(), this.validateUsernames()];

    if (this.value.size === 0) {
      validators.push(this.validateBaseBundle());
    }

    const valid = validators.every(Boolean);

    // set fields as dirty so we force show any errors on submit
    // assuming the user never touched the fields
    runInAction(() => {
      this.dirty.forEach((_, k) => this.dirty.set(k, true));
    });

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
        this.domainSuggestions = mailstack_DomainPurchaseSuggestions.filter(
          (d) => !this.baseBundle.has(d) && !this.extendedBundle.has(d),
        );
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
      this.root.ui.toastError(
        `We're unable to validate your domains`,
        'validate-domains',
      );
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async buyDomains(paymentIntentId: string) {
    try {
      await this.service.buyDomains({
        test: false,
        paymentIntentId,
        domains: [...this.baseBundle, ...this.extendedBundle],
        amount: this.totalAmount,
        username: this.usernames.filter((v) => v !== ''),
        redirectWebsite: this.redirectUrl,
      });
      this.bootstrap();
    } catch (err) {
      this.root.ui.toastError(
        `We couldn't register some of your domains`,
        'buy-domains',
      );
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
