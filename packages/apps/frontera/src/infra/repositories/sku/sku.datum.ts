import { SkusQuery } from './queries/skus.generated';

export type SkuDatum = NonNullable<SkusQuery['skus'][0]>;
