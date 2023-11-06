type RefArray = Array<{
  __ref: string;
}>;
interface Props {
  existing: RefArray;
  incoming: RefArray;
}
export function getUniqueReferenceArray({
  existing,
  incoming,
}: Props): RefArray {
  const merged = existing ? existing.slice(0) : [];
  const existingIds = existing ? existing.map((item) => item.__ref) : [];
  incoming.forEach((item) => {
    if (existingIds.indexOf(item.__ref) < 0) {
      merged.push(item);
    }
  });

  return merged;
}
