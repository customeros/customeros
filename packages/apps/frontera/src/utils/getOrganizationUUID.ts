export const getOrganizationUUID = (path: string): string | null => {
  const match = path.match(
    /\/organization\/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})(?:\/|$)/i,
  );

  return match ? match[1] : null;
};
