import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

export class FilesStore {
  values: Map<string, string> = new Map();
  loaders: Map<string, boolean> = new Map();
  errors: Map<string, string> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
  }

  async download(fileId: string) {
    if (this.values.has(fileId)) return;

    try {
      this.loaders.set(fileId, true);

      const res = await this.transport.http.get(
        `/fs/file/${fileId}/download?inline=true`,
        {
          responseType: 'blob',
        },
      );

      runInAction(() => {
        const url = URL.createObjectURL(res.data);
        this.values.set(fileId, url);
      });
    } catch (err) {
      runInAction(() => {
        this.errors.set(fileId, (err as Error).message);
      });
    } finally {
      runInAction(() => {
        this.loaders.set(fileId, false);
      });
    }
  }

  clear(fileId: string) {
    const url = this.values.get(fileId);
    url && URL.revokeObjectURL(url);

    this.values.delete(fileId);
    this.errors.delete(fileId);
    this.loaders.delete(fileId);
  }

  error(fileId: string) {
    return this.errors.get(fileId);
  }
  isLoading(fileId: string) {
    return this.loaders.get(fileId);
  }
}
