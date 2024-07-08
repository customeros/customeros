import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Contact, Organization } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Download02 } from '@ui/media/icons/Download02.tsx';
import { csvDataMapper as contactCsvDataMapper } from '@organizations/components/Columns/contacts';
import { csvDataMapper as orgCsvDataMapper } from '@organizations/components/Columns/organizations';

const getTableName = (tableViewName: string | undefined) => {
  switch (tableViewName) {
    case 'Targets':
      return 'targets';
    case 'Customers':
      return 'customers';
    case 'Contacts':
      return 'contacts';
    case 'Leads':
      return 'leads';
    case 'Churn':
      return 'churned';
    case 'All orgs':
      return 'organizations';
    default:
      return 'organizations';
  }
};
const convertToCSV = (objArray: Array<Array<string>>): string => {
  return objArray.map((row) => Object.values(row).join(',')).join('\r\n');
};
export const DownloadCsvButton = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableViewName = tableViewDef?.value.name;
  const tableName = getTableName(tableViewName);

  const csvDataMapper =
    tableViewDef?.value.tableType === 'CONTACTS'
      ? contactCsvDataMapper
      : orgCsvDataMapper;

  const handleGetData = (): Array<Array<string>> => {
    const visibleColumns = tableViewDef?.value.columns?.filter(
      (column) => column.visible,
    );
    const headers = visibleColumns?.map((column) =>
      column.columnType.split('_').join(' '),
    ) as Array<string>;

    const data =
      store.ui.filteredTable?.map((row) => {
        return visibleColumns?.map((column) => {
          const mapper: (d: Organization | Contact) => string =
            csvDataMapper?.[column.columnType as keyof typeof csvDataMapper];
          const rowData = (row as ContactStore | OrganizationStore)?.value as
            | Organization
            | Contact;

          return mapper ? mapper?.(rowData) : '';
        }) as Array<string>;
      }) || [];

    return [headers, ...data] as Array<Array<string>>;
  };

  const downloadCSV = () => {
    const data = handleGetData();
    const csvData = new Blob([convertToCSV(data)], { type: 'text/csv' });
    const csvURL = URL.createObjectURL(csvData);
    const link = document.createElement('a');
    link.href = csvURL;
    link.download = `${tableName}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <Tooltip label='Export view as CSV'>
      <IconButton
        aria-label='Download CSV'
        className={'mr-3'}
        size='xs'
        icon={<Download02 />}
        onClick={downloadCSV}
        variant='ghost'
      />
    </Tooltip>
  );
});
