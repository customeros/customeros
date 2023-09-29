import { Grid } from '@ui/layout/Grid';

export default async function AuthPageLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Grid templateColumns={{ base: '1fr', md: '1fr 1fr' }} h='100vh'>
      {children}
    </Grid>
  );
}
