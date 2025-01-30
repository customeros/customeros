import { RootStore } from '@store/root';
import { SkuRepository } from '@infra/repositories/sku';

import { unwrap } from '@utils/unwrap';
import { SkuInput } from '@shared/types/__generated__/graphql.types';

export class SkuService {
  private root = RootStore.getInstance();
  private skuRepository = SkuRepository.getInstance();

  constructor() {}

  public async createSku(data: SkuInput) {
    const tempId = this.root.skus.add(data);

    const [res, err] = await unwrap(
      this.skuRepository.saveSku({
        input: {
          ...data,
        },
      }),
    );

    if (res) {
      this.root.skus.add(res.sku_Save);
      this.root.skus.delete(tempId);
      this.root.ui.toastSuccess('Product created', 'product-add-success');
    }

    if (err) {
      console.error(err);
      this.root.skus.delete(tempId);

      this.root.ui.toastError(
        'Something went wrong creating this product',
        'product-add-failed',
      );

      return [null, err];
    }

    return [res, err];
  }

  public async editSku(data: SkuInput & { id: string }) {
    const prevData = this.root.skus.getById(data.id)?.value;

    this.root.skus.edit(data);

    const [res, err] = await unwrap(
      this.skuRepository.saveSku({
        input: {
          ...data,
        },
      }),
    );

    if (res?.sku_Save) {
      this.root.ui.toastSuccess('Product updated', 'product-edit-success');
    }

    if (err) {
      console.error(err);

      if (prevData) {
        this.root.skus.edit(prevData);
      }

      this.root.ui.toastError(
        'Something went wrong updating this product',
        'product-edit-failed',
      );

      return [null, err];
    }

    return [res, err];
  }

  public async archiveSku(id: string) {
    const deletedSku = this.root.skus.getById(id);

    this.root.skus.delete(id);

    const [res, err] = await unwrap(
      this.skuRepository.archiveSku({
        id,
      }),
    );

    if (res?.sku_Archive?.result) {
      this.root.ui.toastSuccess('Product archived', 'product-archive-success');
    }

    if (err) {
      console.error(err);

      if (deletedSku?.value) {
        this.root.skus.add(deletedSku.value);
      }

      this.root.ui.toastError(
        'Something went wrong archiving this product',
        'product-archive-failed',
      );

      return [null, err];
    }

    return [res, err];
  }
}
