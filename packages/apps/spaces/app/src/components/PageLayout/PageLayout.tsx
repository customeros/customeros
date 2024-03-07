import { Grid } from '@ui/layout/Grid';

interface PageLayoutProps {
  children: React.ReactNode;
}

export const PageLayout = ({ children }: PageLayoutProps) => {
  return (
    <Grid
      h='100vh'
      bg='gray.25'
      templateColumns='200px 1fr'
      templateRows='1fr'
      transition='all ease 0.25s'
      templateAreas={`"sidebar content"`}
    >
      {children}
    </Grid>
  );
};
