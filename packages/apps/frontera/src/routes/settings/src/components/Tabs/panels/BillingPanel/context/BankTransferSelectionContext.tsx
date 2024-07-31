import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { BankAccount } from '@graphql/types';

interface BankTransferSelectionContextMethods {
  hoveredAccount: null | BankAccount;
  focusedAccount: null | BankAccount;
  defaultSelectedAccount: null | BankAccount;
  setAccounts: (accounts: Array<BankAccount>) => void;
  setHoverAccount: (account: BankAccount | null) => void;
  setFocusAccount: (account: BankAccount | null) => void;
}

const BankTransferSelectionContext =
  createContext<BankTransferSelectionContextMethods>({
    defaultSelectedAccount: null,
    hoveredAccount: null,
    focusedAccount: null,
    setHoverAccount: () => null,
    setFocusAccount: () => null,
    setAccounts: () => null,
  });

export const useBankTransferSelectionContext = () => {
  return useContext(BankTransferSelectionContext);
};

export const BankTransferSelectionContextProvider = ({
  children,
}: PropsWithChildren) => {
  const [accounts, setAccounts] = useState<Array<BankAccount>>([]);
  const [defaultSelectedAccount, setDefaultSelectedAccount] =
    useState<BankAccount | null>(null);
  const [hoveredAccount, setHoveredAccount] = useState<BankAccount | null>(
    null,
  );
  const [focusedAccount, setFocusedAccount] = useState<BankAccount | null>(
    null,
  );

  const handleFocusAccount = (account: BankAccount | null) => {
    if (account) {
      setFocusedAccount(account);

      return;
    }
    setFocusedAccount(null);
  };

  const handleHoverAccount = (account: BankAccount | null) => {
    if (account) {
      setHoveredAccount(account);

      return;
    }
    setHoveredAccount(null);
  };

  useEffect(() => {
    setDefaultSelectedAccount(accounts[0]);
  }, [accounts]);

  return (
    <BankTransferSelectionContext.Provider
      value={{
        defaultSelectedAccount: defaultSelectedAccount,
        hoveredAccount,
        focusedAccount,
        setHoverAccount: handleHoverAccount,
        setFocusAccount: handleFocusAccount,
        setAccounts,
      }}
    >
      {children}
    </BankTransferSelectionContext.Provider>
  );
};
