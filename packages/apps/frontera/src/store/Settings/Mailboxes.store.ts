import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import { override, runInAction, makeObservable } from 'mobx';

import { RegisteredBuyDomainWithMailboxes } from '@shared/types/__generated__/graphql.types';

import { MailboxStore } from './Mailbox.store';
import { MailboxesService } from './__service__/Mailboxes/Mailboxes.service';
import { CheckUnavailableDomainsQueryVariables } from './__service__/Mailboxes/getMailstackCheckUnavailableDomains.generated';

export class MailboxesStore extends SyncableGroup<
  RegisteredBuyDomainWithMailboxes,
  MailboxStore
> {
  private service: MailboxesService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, MailboxStore);
    this.service = MailboxesService.getInstance(transport);

    makeObservable<MailboxesStore>(this, {
      channelName: override,
    });
  }

  async bootstrap() {
    try {
      const { mailstack_RegisteredBuyDomainsWithMailboxes } =
        await this.service.getRegisteredMailboxes();

      const mailboxes = mailstack_RegisteredBuyDomainsWithMailboxes.map(
        (template) =>
          ({
            createdAt: template.createdAt,
            domain: {
              domain: template.domain.domain,
              status: template.domain.status,
            },
            id: template.id,
            mailboxes: template.mailboxes.map((mailbox) => ({
              mailbox: mailbox.mailbox,
              status: mailbox.status,
            })),
          } as RegisteredBuyDomainWithMailboxes),
      );

      this.load(mailboxes as RegisteredBuyDomainWithMailboxes[], {
        getId: (data: RegisteredBuyDomainWithMailboxes) => data.id,
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

  async getDomainsSuggestions(payload: string) {
    try {
      const reponse = await this.service.getMailstackDomainsSuggestions({
        domain: payload,
      });

      return reponse.mailstack_DomainPurchaseSuggestions;
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

  async getMailstackCheckUnvalidDomains(
    payload: CheckUnavailableDomainsQueryVariables,
  ) {
    try {
      const reponse = await this.service.getMailstackCheckUnavailableDomains({
        domains: payload.domains,
      });

      return reponse.mailstack_CheckUnavailableDomains;
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

  // async save(payload: CustomFieldTemplateInput) {
  //   try {
  //     const response = await this.service.saveCustomField({
  //       validValues: payload.validValues,
  //       name: payload.name,
  //       type: payload.type,
  //       entityType: payload.entityType,
  //     });

  //     const customField: CustomFieldTemplateInput = {
  //       id: response.customFieldTemplate_Save.id,
  //       name: payload.name,
  //       entityType: payload.entityType,
  //       type: payload?.type,
  //       validValues: payload.validValues,
  //     };

  //     runInAction(() => {
  //       this.load(
  //         [
  //           {
  //             id: customField.id || '',
  //             name: customField.name || '',
  //             value: '',
  //             createdAt: new Date(),
  //             updatedAt: new Date(),
  //             source: DataSource.Openline,
  //             datatype: CustomFieldDataType.Text,
  //             template: {
  //               name: customField.name || '',
  //               id: customField.id || '',
  //               validValues: customField.validValues || [],
  //               createdAt: new Date(),
  //               updatedAt: new Date(),
  //               entityType: customField.entityType || EntityType.Organization,
  //               type: customField.type || CustomFieldTemplateType.FreeText,
  //             },
  //           },
  //         ],
  //         {
  //           getId: (data: CustomField) => data.id,
  //         },
  //       );
  //     });
  //   } catch (err) {
  //     runInAction(() => {
  //       this.error = (err as Error).message;
  //     });
  //   }
  // }

  // async deleteCustomField(id: string) {
  //   const customField = this.value.get(id);

  //   try {
  //     await this.service.deleteCustomField(id);

  //     if (customField) {
  //       runInAction(() => {
  //         this.value.delete(id);
  //       });
  //     }
  //   } catch (err) {
  //     runInAction(() => {
  //       this.error = (err as Error).message;
  //     });
  //   }
  // }

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
