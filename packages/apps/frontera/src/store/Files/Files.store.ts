import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import { toastError } from '@ui/presentation/Toast';

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
  private getFileExtension(mimeType: string): string {
    const mimeToExtension: { [key: string]: string } = {
      'application/pdf': 'pdf',
      'image/jpeg': 'jpg',
      'image/png': 'png',
      'text/plain': 'txt',
      // invoices come as octet-stream but we want to use pdf extension
      'application/octet-stream': 'pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document':
        'docx',
      'image/gif': 'gif',
      'image/webp': 'webp',
      'image/svg+xml': 'svg',
      'image/tiff': 'tiff',
    };

    return mimeToExtension[mimeType] || 'unknown';
  }

  async downloadAttachment(fileId: string, fileName: string) {
    if (this.values.has(fileId)) return;

    try {
      this.loaders.set(fileId, true);

      const res = await this.transport.http.get(`/fs/file/${fileId}/download`, {
        responseType: 'blob',
      });
      const mimeType = res.data.type;
      const blobUrl = window.URL.createObjectURL(res.data);
      const fileExtension = this.getFileExtension(mimeType);

      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = `${fileName}.${fileExtension}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      setTimeout(() => {
        window.URL.revokeObjectURL(blobUrl);
      }, 100);
    } catch (err) {
      runInAction(() => {
        this.errors.set(fileId, (err as Error).message);
        toastError(
          'Something went wrong while downloading the file',
          'download-attachment-error',
        );
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
