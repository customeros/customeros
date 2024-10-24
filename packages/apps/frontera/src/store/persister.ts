import LocalForage from 'localforage';

export type PersisterInstance = LocalForage;

export class Persister {
  static DB_NAME = 'customerDB';
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
}
