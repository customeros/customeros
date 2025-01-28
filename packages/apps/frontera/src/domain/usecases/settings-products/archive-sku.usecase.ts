import { action } from 'mobx';
import { SkuService } from '@domain/services';

export class ArchiveSkuUsecase {
  private service = new SkuService();

  constructor() {}

  @action
  archiveSku(id: string) {
    this.service.archiveSku(id);
  }
}
