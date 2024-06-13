import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { TableIdType, type TableViewDef } from '@graphql/types';

import { TableViewDefStore } from './TableViewDef.store';

export class TableViewDefsStore implements GroupStore<TableViewDef> {
  value: Map<string, TableViewDefStore> = new Map();
  isLoading = false;
  channel?: Channel;
  version: number = 0;
  history: GroupOperation[] = [];
  isBootstrapped = false;
  error: string | null = null;
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<TableViewDef>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncableGroup(this, {
      channelName: 'TableViewDefs',
      ItemStore: TableViewDefStore,
      getItemId: (item) => item.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    if (this.root.demoMode) {
      this.load(mock.data.tableViewDefs as TableViewDef[]);
      this.isBootstrapped = true;

      return;
    }

    try {
      this.isLoading = true;
      const res =
        await this.transport.graphql.request<TABLE_VIEW_DEFS_QUERY_RESULT>(
          TABLE_VIEW_DEFS_QUERY,
        );

      this.load(res?.tableViewDefs);
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

  getById(id: string) {
    return this.value.get(id);
  }

  toArray(): TableViewDefStore[] {
    return Array.from(this.value)?.flatMap(
      ([, tableViewDefStore]) => tableViewDefStore,
    );
  }

  get defaultPreset() {
    const preset = this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Organizations,
    )?.value.id;

    return preset;
  }
}

type TABLE_VIEW_DEFS_QUERY_RESULT = { tableViewDefs: TableViewDef[] };
const TABLE_VIEW_DEFS_QUERY = gql`
  query tableViewDefs {
    tableViewDefs {
      id
      name
      tableType
      tableId
      order
      icon
      filters
      sorting
      columns {
        columnType
        width
        visible
      }
      createdAt
      updatedAt
    }
  }
`;

const mock = {
  data: {
    tableViewDefs: [
      {
        id: '1052',
        name: 'Customers',
        tableType: 'ORGANIZATIONS',
        tableId: 'CUSTOMERS',
        order: 1,
        icon: 'CheckHeart',
        filters:
          '{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["CUSTOMER"]}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'ORGANIZATIONS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_WEBSITE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RELATIONSHIP',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RENEWAL_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_ONBOARDING_STATUS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_FORECAST_ARR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_OWNER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:34.933188Z',
        updatedAt: '2024-06-07T05:11:34.933188Z',
      },
      {
        id: '1046',
        name: 'Monthly renewals',
        tableType: 'RENEWALS',
        tableId: 'MONTHLY_RENEWALS',
        order: 1,
        icon: 'ClockFastForward',
        filters:
          '{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"MONTHLY","operation":"EQ","includeEmpty":false}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'RENEWALS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_RENEWAL_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_FORECAST_ARR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_RENEWAL_LIKELIHOOD',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_OWNER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:24.826817Z',
        updatedAt: '2024-06-07T05:11:24.826817Z',
      },
      {
        id: '1047',
        name: 'Quarterly renewals',
        tableType: 'RENEWALS',
        tableId: 'QUARTERLY_RENEWALS',
        order: 2,
        icon: 'ClockFastForward',
        filters:
          '{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"QUARTERLY","operation":"EQ","includeEmpty":false}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'RENEWALS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_RENEWAL_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_FORECAST_ARR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_RENEWAL_LIKELIHOOD',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_OWNER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:29.86154Z',
        updatedAt: '2024-06-07T05:11:29.86154Z',
      },
      {
        id: '1054',
        name: 'Leads',
        tableType: 'ORGANIZATIONS',
        tableId: 'LEADS',
        order: 3,
        icon: 'SwitchHorizontal01',
        filters:
          '{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["LEAD"]}},{"filter":{"property":"ORGANIZATIONS_CREATED_DATE","value":"1997-01-24","active":false,"caseSensitive":false,"includeEmpty":false,"operation":"LTE"}},{"filter":{"property":"ORGANIZATIONS_LEAD_SOURCE","value":[],"active":false,"caseSensitive":false,"includeEmpty":false,"operation":"IN"}},{"filter":{"property":"ORGANIZATIONS_YEAR_FOUNDED","value":["Computers Electronics and Technology"],"active":true,"caseSensitive":false,"includeEmpty":false,"operation":"IN"}},{"filter":{"property":"ORGANIZATIONS_INDUSTRY","value":[],"active":false,"caseSensitive":false,"includeEmpty":false,"operation":"IN"}},{"filter":{"property":"ORGANIZATIONS_NAME","value":"","active":false,"caseSensitive":false,"includeEmpty":false,"operation":"CONTAINS"}},{"filter":{"property":"ORGANIZATIONS_WEBSITE","value":"","active":false,"caseSensitive":false,"includeEmpty":false,"operation":"CONTAINS"}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'ORGANIZATIONS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_WEBSITE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_SOCIALS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_CREATED_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LEAD_SOURCE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_EMPLOYEE_COUNT',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_YEAR_FOUNDED',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_INDUSTRY',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:39.960336Z',
        updatedAt: '2024-06-12T12:04:20.855631Z',
      },
      {
        id: '1048',
        name: 'Annual renewals',
        tableType: 'RENEWALS',
        tableId: 'QUARTERLY_RENEWALS',
        order: 3,
        icon: 'ClockFastForward',
        filters:
          '{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"ANNUALLY","operation":"EQ","includeEmpty":false}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'RENEWALS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_RENEWAL_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_FORECAST_ARR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_RENEWAL_LIKELIHOOD',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_OWNER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'RENEWALS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:29.87764Z',
        updatedAt: '2024-06-07T05:11:29.87764Z',
      },
      {
        id: '1055',
        name: 'Nurture',
        tableType: 'ORGANIZATIONS',
        tableId: 'NURTURE',
        order: 4,
        icon: 'HeartHand',
        filters:
          '{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["TARGET"]}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["PROSPECT"]}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'ORGANIZATIONS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_WEBSITE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_SOCIALS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:39.97271Z',
        updatedAt: '2024-06-07T05:11:39.97271Z',
      },
      {
        id: '1049',
        name: 'Upcoming',
        tableType: 'INVOICES',
        tableId: 'UPCOMING_INVOICES',
        order: 4,
        icon: 'InvoiceUpcoming',
        filters:
          '{"AND":[{"filter":{"property":"INVOICE_PREVIEW","value":true}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'INVOICES_INVOICE_PREVIEW',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_CONTRACT',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_BILLING_CYCLE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_ISSUE_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_DUE_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_AMOUNT',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_INVOICE_STATUS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_ISSUE_DATE_PAST',
            width: 100,
            visible: false,
          },
          {
            columnType: 'INVOICES_PAYMENT_STATUS',
            width: 100,
            visible: false,
          },
        ],
        createdAt: '2024-06-07T05:11:29.890433Z',
        updatedAt: '2024-06-07T05:11:29.890433Z',
      },
      {
        id: '1050',
        name: 'Past',
        tableType: 'INVOICES',
        tableId: 'PAST_INVOICES',
        order: 5,
        icon: 'InvoiceCheck',
        filters:
          '{"AND":[{"filter":{"property":"INVOICE_DRY_RUN","value":false}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'INVOICES_INVOICE_NUMBER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_CONTRACT',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_BILLING_CYCLE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_ISSUE_DATE_PAST',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_DUE_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_AMOUNT',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_PAYMENT_STATUS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'INVOICES_ISSUE_DATE',
            width: 100,
            visible: false,
          },
          {
            columnType: 'INVOICES_INVOICE_STATUS',
            width: 100,
            visible: false,
          },
        ],
        createdAt: '2024-06-07T05:11:29.903837Z',
        updatedAt: '2024-06-07T05:11:29.903837Z',
      },
      {
        id: '1051',
        name: 'All orgs',
        tableType: 'ORGANIZATIONS',
        tableId: 'ORGANIZATIONS',
        order: 5,
        icon: 'Building07',
        filters:
          '{"AND":[{"filter":{"property":"ORGANIZATIONS_LAST_TOUCHPOINT","value":{"after":"2024-06-04","types":[]},"active":false,"caseSensitive":false,"includeEmpty":false,"operation":"EQ"}},{"filter":{"property":"ORGANIZATIONS_OWNER","value":["3f30a3e9-18a0-4a27-83f8-5340941d9e37"],"active":false,"caseSensitive":false,"includeEmpty":false,"operation":"IN"}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'ORGANIZATIONS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_WEBSITE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RELATIONSHIP',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RENEWAL_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_ONBOARDING_STATUS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_FORECAST_ARR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_OWNER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:34.917558Z',
        updatedAt: '2024-06-12T07:50:57.145295Z',
      },
      {
        id: '1108',
        name: 'Churn',
        tableType: 'ORGANIZATIONS',
        tableId: 'CHURN',
        order: 5,
        icon: 'BrokenHeart',
        filters:
          '{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["FORMER_CUSTOMER"]}},{"filter":{"property":"ORGANIZATIONS_LTV","value":[0,1253334],"active":false,"caseSensitive":false,"includeEmpty":false,"operation":"CONTAINS"}},{"filter":{"property":"ORGANIZATIONS_NAME","value":"","active":false,"caseSensitive":false,"includeEmpty":false,"operation":"CONTAINS"}}]}',
        sorting: '{"id": "ORGANIZATIONS_CHURN_DATE", "desc": true}',
        columns: [
          {
            columnType: 'ORGANIZATIONS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_CHURN_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LTV',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T13:17:12.08309Z',
        updatedAt: '2024-06-12T07:51:01.62867Z',
      },
      {
        id: '1053',
        name: 'My portfolio',
        tableType: 'ORGANIZATIONS',
        tableId: 'MY_PORTFOLIO',
        order: 6,
        icon: 'Briefcase01',
        filters:
          '{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"OWNER_ID","value":["6503cb51-21c8-451c-b956-1dc28c72545a"]}}]}',
        sorting: '',
        columns: [
          {
            columnType: 'ORGANIZATIONS_AVATAR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_NAME',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_WEBSITE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RELATIONSHIP',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_RENEWAL_DATE',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_ONBOARDING_STATUS',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_FORECAST_ARR',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_OWNER',
            width: 100,
            visible: true,
          },
          {
            columnType: 'ORGANIZATIONS_LAST_TOUCHPOINT',
            width: 100,
            visible: true,
          },
        ],
        createdAt: '2024-06-07T05:11:34.946397Z',
        updatedAt: '2024-06-07T05:11:34.946397Z',
      },
    ],
  },
};
