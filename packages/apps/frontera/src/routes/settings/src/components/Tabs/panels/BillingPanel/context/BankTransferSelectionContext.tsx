import { useState, useContext, createContext, PropsWithChildren } from 'react';

interface BankTransferSelectionContextMethods {
  hoveredAccount: null | string;
  focusedAccount: null | string;
  setHoverAccount: (account: string | null) => void;
  setFocusAccount: (account: string | null) => void;
}

const BankTransferSelectionContext =
  createContext<BankTransferSelectionContextMethods>({
    hoveredAccount: null,
    focusedAccount: null,
    setHoverAccount: () => null,
    setFocusAccount: () => null,
  });

export const useBankTransferSelectionContext = () => {
  return useContext(BankTransferSelectionContext);
};

export const BankTransferSelectionContextProvider = ({
  children,
}: PropsWithChildren) => {
  const [hoveredAccount, setHoveredAccount] = useState<string | null>(null);
  const [focusedAccount, setFocusedAccount] = useState<string | null>(null);

  const handleFocusAccount = (account: string | null) => {
    if (account) {
      setFocusedAccount(account);

      return;
    }
    setFocusedAccount(null);
  };

  const handleHoverAccount = (account: string | null) => {
    if (account) {
      setHoveredAccount(account);

      return;
    }
    setHoveredAccount(null);
  };

  return (
    <BankTransferSelectionContext.Provider
      value={{
        hoveredAccount,
        focusedAccount,
        setHoverAccount: handleHoverAccount,
        setFocusAccount: handleFocusAccount,
      }}
    >
      {children}
    </BankTransferSelectionContext.Provider>
  );
};
