import { Grid } from '@ui/layout/Grid';

interface PageLayoutProps {
  children: React.ReactNode;
}

export const PageLayout = ({ children }: PageLayoutProps) => {
  return (
    <Grid
      p='2'
      gap='2'
      h='100vh'
      templateColumns='200px 1fr'
      transition='all ease 0.25s'
      templateAreas={`"sidebar content"`}
      bgGradient='linear(to-t, #EAECF0, #F3F4F7)'
    >
      {children}
    </Grid>
  );
};
