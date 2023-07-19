'use client';
import { Sidebar } from '../Sidebar/Sidebar';

import { Grid, GridItem } from '@ui/layout/Grid';

interface PageLayoutProps {
  isOwner: boolean;
}

export const PageLayout = ({
  isOwner,
  children,
}: React.PropsWithChildren<PageLayoutProps>) => {
  return (
    <Grid
      p='2'
      gap='2'
      h='100vh'
      templateColumns='200px 1fr'
      transition='all ease 0.25s'
      templateAreas={`"sidebar content"`}
      bgGradient='linear(to-t, gray.200, gray.50)'
    >
      <Sidebar isOwner={isOwner} />
      <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
        {children}
      </GridItem>
    </Grid>
  );
};
