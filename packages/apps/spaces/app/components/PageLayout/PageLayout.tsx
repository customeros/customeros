import { Sidebar } from '../Sidebar/Sidebar';

import styles from './PageLayout.module.scss';

interface PageLayoutProps {
  isOwner: boolean;
}

export const PageLayout = ({
  isOwner,
  children,
}: React.PropsWithChildren<PageLayoutProps>) => {
  return (
    <div className={styles.pageContent}>
      <Sidebar isOwner={isOwner} />
      <div className={styles.innerWrapper}>{children}</div>
    </div>
  );
};
