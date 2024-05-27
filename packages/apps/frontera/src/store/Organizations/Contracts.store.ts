import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Contract } from '@graphql/types';

import { ContractStore } from './Contract.store';

export class ContractsStore implements GroupStore<Contract> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<Contract>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Contract>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: `Contracts:${this.root.session.value.tenant}`,
      getItemId: (item) => item?.metadata?.id,
      ItemStore: ContractStore,
    });
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;
      // const res = await this.transport.graphql.request(``, {
      //   pagination: { limit: 1000, page: 0 },
      // });

      // this.load(dashboardView_Organizations.content);
      runInAction(() => {
        this.isBootstrapped = true;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

const _CONTRACTS_QUERY = gql`
  query getContracts($id: ID!) {
    organization(id: $id) {
      id
      name
      note
      accountDetails {
        renewalSummary {
          arrForecast
          maxArrForecast
          renewalLikelihood
        }
      }
      contracts {
        metadata {
          id
          created
          source
          lastUpdated
        }
        contractName
        serviceStarted
        contractSigned
        contractEnded
        contractStatus
        committedPeriodInMonths
        approved

        contractUrl
        billingCycle
        billingEnabled
        currency
        invoiceEmail
        autoRenew

        billingDetails {
          nextInvoicing
          postalCode
          country
          locality
          addressLine1
          addressLine2
          invoiceNote
          organizationLegalName
          billingCycle
          invoicingStarted
          region
          dueDays
          billingEmail
          billingEmailCC
          billingEmailBCC
        }
        upcomingInvoices {
          metadata {
            id
          }
          invoicePeriodEnd
          invoicePeriodStart
          status
          issued
          amountDue
          due
          currency
          invoiceLineItems {
            metadata {
              id
              created
            }

            quantity
            subtotal
            taxDue
            total
            price
            description
          }
          contract {
            billingDetails {
              canPayWithBankTransfer
            }
          }
          status
          invoiceNumber
          invoicePeriodStart
          invoicePeriodEnd
          invoiceUrl
          due
          issued
          subtotal
          taxDue
          currency
          note
          customer {
            name
            email
            addressLine1
            addressLine2
            addressZip
            addressLocality
            addressCountry
            addressRegion
          }
          provider {
            name
            addressLine1
            addressLine2
            addressZip
            addressLocality
            addressCountry
          }
        }
        opportunities {
          id
          comments
          internalStage
          internalType
          amount
          maxAmount
          name
          renewalLikelihood
          renewalAdjustedRate
          renewalUpdatedByUserId
          renewedAt
          updatedAt

          owner {
            id
            firstName
            lastName
            name
          }
        }
        contractLineItems {
          metadata {
            id
            created
            lastUpdated
            source
            appSource
            sourceOfTruth
          }
          description
          billingCycle
          price
          quantity
          comments
          serviceEnded
          parentId
          serviceStarted
          tax {
            salesTax
            vat
            taxRate
          }
        }
      }
    }
  }
`;
