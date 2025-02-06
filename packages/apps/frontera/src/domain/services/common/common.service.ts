import { RootStore } from '@store/root';

export class CommonService {
  private root = RootStore.getInstance();

  constructor() {}

  public async switchWorkspace(tenant: string) {
    const r = await this.root.session.transport.http.get<{
      redirectUrl: string;
    }>('/switchWorkspace?tenant=' + tenant);

    await this.root.session.clearSession();
    window.location.href = r.data.redirectUrl;

    return 'success';
  }
}
