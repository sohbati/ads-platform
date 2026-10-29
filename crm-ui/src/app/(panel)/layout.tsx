import { AppShell } from "@/components/layout/AppShell";
import { SessionProvider } from "@/lib/auth/SessionProvider";
import { requireSession } from "@/lib/auth/guards";

export default async function PanelLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = await requireSession();
  return (
    <SessionProvider session={session}>
      <AppShell>{children}</AppShell>
    </SessionProvider>
  );
}

