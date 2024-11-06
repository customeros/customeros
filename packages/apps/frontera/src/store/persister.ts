import LocalForage from 'localforage';

export type PersisterInstance = LocalForage;

export class Persister {
  static DB_NAME = 'customerDB';
  private static version = 2.9;
  private static instances: Map<string, PersisterInstance> = new Map();
  private static sharedInstances: Map<string, PersisterInstance> = new Map();

  constructor() {}

  public static getInstance(key: string) {
    if (!Persister.instances.has(key)) {
      const newInstance = LocalForage.createInstance({
        driver: LocalForage.INDEXEDDB,
        name: Persister.DB_NAME,
        storeName: key,
      });

      Persister.instances.set(key, newInstance);
    }

    return Persister.instances.get(key);
  }

  public static getSharedInstance(key: string) {
    if (!Persister.sharedInstances.has(key)) {
      const newInstance = LocalForage.createInstance({
        driver: LocalForage.INDEXEDDB,
        name: 'customerDB_shared',
        storeName: key,
      });

      Persister.sharedInstances.set(key, newInstance);
    }

    return Persister.sharedInstances.get(key);
  }

  public static setTenant(tenant: string) {
    Persister.DB_NAME = `customerDB_${tenant}`;
  }

  public static async attemptPurge(opt?: { force?: boolean }) {
    try {
      const instance = Persister.getSharedInstance('Meta');
      const version = await instance?.getItem('version');

      if (!version) {
        await instance?.setItem('version', Persister.version);

        return;
      }

      if (opt?.force || version !== Persister.version) {
        const dbs = await indexedDB.databases();
        const dbNames = dbs
          .map((db) => db.name)
          .filter(
            (n) => n?.startsWith('customerDB_') && !n?.includes('shared'),
          );

        const drops = dbNames.map((name) => LocalForage.dropInstance({ name }));

        await Promise.all(drops);
        await instance?.setItem('version', Persister.version);
      }
    } catch (err) {
      console.error('Failed to attempt purging indexedDB', err);
    }
  }
}
