import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Download02 } from '@ui/media/icons/Download02.tsx';
import { csvDataMapper } from '@organizations/components/Columns/organizations';

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
const convertToCSV = (objArray: Array<Record<string, string>>) => {
  return objArray.map((row) => Object.values(row).join(',')).join('\r\n');
};
export const DownloadCsvButton = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableViewName = tableViewDef?.value.name;
  const tableName = getTableName(tableViewName);

  const handleGetData = () => {
    const headers = tableViewDef?.value.columns?.map((column) =>
      column.columnType.split('_').join(' '),
    );

    const data = store?.ui?.filteredTable?.map((row) => {
      const visibleColumns = tableViewDef?.value.columns?.filter(
        (column) => column.visible,
      );

      return visibleColumns?.map((column) =>
        csvDataMapper?.[column.columnType]?.(row?.value),
      );
    });

    return [headers, ...data];
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
        icon={<Download02 />}
        onClick={downloadCSV}
        variant='ghost'
      />
    </Tooltip>
  );
});
