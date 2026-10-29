"use client";

import { createContext, useContext, useMemo } from "react";
import type { AdminSession } from "./session.types";

const SessionContext = createContext<AdminSession | null>(null);

export function SessionProvider({
  session,
  children,
}: {
  session: AdminSession;
  children: React.ReactNode;
}) {
  const value = useMemo(() => session, [session]);
  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession(): AdminSession {
  const session = useContext(SessionContext);
  if (!session) {
    throw new Error("useSession must be used within SessionProvider");
  }
  return session;
}
