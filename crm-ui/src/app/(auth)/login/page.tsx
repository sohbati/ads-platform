import { redirect } from "next/navigation";
import { LoginForm } from "@/features/auth/components/LoginForm";
import { getSession } from "@/lib/auth/session";

export default async function LoginPage() {
  let signedIn = false;
  try {
    signedIn = (await getSession()) !== null;
  } catch {
    // Backend unreachable: show the form; signing in will report the error.
  }
  if (signedIn) {
    redirect("/dashboard");
  }
  return (
    <main className="login-shell">
      <LoginForm />
    </main>
  );
}
