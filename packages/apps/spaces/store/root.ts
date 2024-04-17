import { TransportLayer } from './transport';
import { TableViewDefsStore } from './TableViewDefs/TableViewDefs.store';

export class RootStore {
  tableViewDefsStore: TableViewDefsStore;

  constructor(private transportLayer: TransportLayer) {
    this.tableViewDefsStore = new TableViewDefsStore(this, transportLayer);
  }
}
