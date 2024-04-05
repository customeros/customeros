import { Search } from './src/components/Search';
import { Preview } from './src/components/Preview';
import { InvoicesTable } from './src/components/InvoicesTable';

export default function InvoicesPage() {
  return (
    <>
      <Search />
      <InvoicesTable />
      <Preview />
    </>
  );
}
