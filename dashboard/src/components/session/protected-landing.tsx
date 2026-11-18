import { redirect } from "next/navigation";
import { SessionLanding } from "@/components/session/session-landing";
import { getServerSession } from "@/lib/server-api";

export async function ProtectedLanding({ destination, title, description }: { destination: "/admin" | "/familia" | "/acceso-pendiente"; title: string; description: string }) {
  const session = await getServerSession();
  if (!session) redirect("/");
  if (session.redirect !== destination) redirect(session.redirect);
  return <SessionLanding user={session.user} title={title} description={description} />;
}
