"use client";

import {
  Activity,
  Building2,
  Clipboard,
  KeyRound,
  LayoutDashboard,
  LogOut,
  Plus,
  Search,
  ShieldCheck,
  UserPlus,
  Users,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { FormEvent, useMemo, useState } from "react";
import {
  APIError,
  createInvitation,
  createTenant,
  getPlatformSnapshot,
  logout,
  revokeUserSessions,
  setMembershipRole,
  setTenantStatus,
  setUserStatus,
  type AuthenticatedUser,
  type AuditEvent,
  type Membership,
  type PlatformSnapshot,
  type TenantRecord,
} from "@/lib/api";

type View = "overview" | "tenants" | "users" | "audit";
type Panel = "tenant" | "invitation" | null;

const statusLabels: Record<string, string> = {
  active: "Activo",
  invited: "Pendiente",
  suspended: "Suspendido",
  deactivated: "Desactivado",
};

const auditLabels: Record<string, string> = {
  "auth.login.succeeded": "Inicio de sesión",
  "auth.logout": "Cierre de sesión",
  "platform.tenant.created": "Organización creada",
  "platform.tenant.status_changed": "Estado de organización actualizado",
  "platform.user.invited": "Usuario invitado",
  "platform.user.status_changed": "Estado de usuario actualizado",
  "platform.user.sessions_revoked": "Sesiones revocadas",
  "platform.membership.role_changed": "Rol actualizado",
  "identity.invitation.accepted": "Invitación aceptada",
  "identity.platform_admin.created": "Administrador de plataforma creado",
};

export function PlatformDashboard({ user, initialSnapshot }: { user: AuthenticatedUser; initialSnapshot: PlatformSnapshot }) {
  const router = useRouter();
  const [snapshot, setSnapshot] = useState(initialSnapshot);
  const [view, setView] = useState<View>("overview");
  const [panel, setPanel] = useState<Panel>(null);
  const [query, setQuery] = useState("");
  const [busy, setBusy] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; text: string } | null>(null);
  const [activationLink, setActivationLink] = useState<string | null>(null);

  const filteredTenants = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    if (!normalized) return snapshot.tenants;
    return snapshot.tenants.filter((tenant) => [tenant.name, tenant.slug, tenant.ownerName, tenant.ownerEmail].some((value) => value?.toLowerCase().includes(normalized)));
  }, [query, snapshot.tenants]);

  const filteredUsers = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    if (!normalized) return snapshot.users;
    return snapshot.users.filter((entry) => [entry.displayName, entry.email, ...entry.memberships.map((item) => item.tenantName)].some((value) => value.toLowerCase().includes(normalized)));
  }, [query, snapshot.users]);

  async function refresh() {
    setSnapshot(await getPlatformSnapshot());
    router.refresh();
  }

  async function runMutation(key: string, action: () => Promise<unknown>, success: string) {
    setBusy(key);
    setFeedback(null);
    try {
      await action();
      await refresh();
      setFeedback({ kind: "success", text: success });
      return true;
    } catch (error) {
      setFeedback({ kind: "error", text: error instanceof APIError ? error.message : "No pudimos completar la acción." });
      return false;
    } finally {
      setBusy(null);
    }
  }

  async function handleCreateTenant(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    const created = await runMutation("create-tenant", () => createTenant({ name: String(data.get("name")), slug: String(data.get("slug")) }), "Organización creada y registrada en auditoría.");
    if (created) form.reset();
  }

  async function handleInvite(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    setBusy("invite-user");
    setFeedback(null);
    setActivationLink(null);
    try {
      const result = await createInvitation({
        tenantId: String(data.get("tenantId")),
        displayName: String(data.get("displayName")),
        email: String(data.get("email")),
        role: String(data.get("role")) as Membership["role"],
      });
      if (result.activationToken) {
        setActivationLink(`${window.location.origin}/activar?token=${encodeURIComponent(result.activationToken)}`);
      }
      await refresh();
      form.reset();
      setFeedback({ kind: "success", text: result.requiresActivation ? "Invitación creada. Copia el enlace; solo se mostrará esta vez." : "El usuario existente fue agregado a la organización." });
    } catch (error) {
      setFeedback({ kind: "error", text: error instanceof APIError ? error.message : "No pudimos crear la invitación." });
    } finally {
      setBusy(null);
    }
  }

  async function handleLogout() {
    setBusy("logout");
    try {
      await logout();
      router.replace("/");
      router.refresh();
    } catch (error) {
      setFeedback({ kind: "error", text: error instanceof APIError ? error.message : "No se pudo cerrar la sesión. Intenta nuevamente." });
      setBusy(null);
    }
  }

  async function copyActivationLink() {
    if (!activationLink) return;
    await navigator.clipboard.writeText(activationLink);
    setFeedback({ kind: "success", text: "Enlace de activación copiado." });
  }

  const metrics = [
    [snapshot.summary.tenantsTotal, "Organizaciones"],
    [snapshot.summary.usersActive, "Usuarios activos"],
    [snapshot.summary.pendingInvitations, "Invitaciones"],
    [snapshot.summary.usersSuspended, "Suspendidos"],
  ] as const;

  return (
    <div className="platform-shell">
      <aside className="platform-rail">
        <div className="platform-brand"><span>APU</span><i aria-hidden="true" /></div>
        <nav aria-label="Navegación de plataforma">
          <NavButton active={view === "overview"} onClick={() => setView("overview")} icon={<LayoutDashboard />} label="Panorama" />
          <NavButton active={view === "tenants"} onClick={() => setView("tenants")} icon={<Building2 />} label="Organizaciones" />
          <NavButton active={view === "users"} onClick={() => setView("users")} icon={<Users />} label="Usuarios" />
          <NavButton active={view === "audit"} onClick={() => setView("audit")} icon={<Activity />} label="Auditoría" />
        </nav>
        <div className="rail-account">
          <span>{user.displayName}</span>
          <small>Administrador de plataforma</small>
          <button type="button" onClick={handleLogout} disabled={busy === "logout"}><LogOut aria-hidden="true" />{busy === "logout" ? "Cerrando…" : "Cerrar sesión"}</button>
        </div>
      </aside>

      <main className="platform-main">
        <header className="platform-header">
          <div><h1>{viewTitle(view)}</h1><p>Control operativo de clientes, accesos y registros de APU.</p></div>
          <div className="platform-coordinate" aria-label="Contexto de operación"><span>LPB / PLATAFORMA</span><strong>{new Intl.DateTimeFormat("es-BO", { dateStyle: "medium", timeZone: "America/La_Paz" }).format(new Date())}</strong></div>
        </header>

        <section className="metric-strip" aria-label="Resumen de plataforma">
          {metrics.map(([value, label]) => <div key={label}><strong>{value}</strong><span>{label}</span></div>)}
        </section>

        {feedback && <div className={`platform-feedback is-${feedback.kind}`} role="status"><span className="registration-point" aria-hidden="true" />{feedback.text}</div>}
        {activationLink && <div className="activation-ticket"><div><strong>Enlace de activación de un solo uso</strong><code>{activationLink}</code></div><button type="button" onClick={copyActivationLink}><Clipboard aria-hidden="true" />Copiar enlace</button></div>}

        <section className={`platform-toolbar${view === "audit" ? " is-audit" : ""}`} aria-label="Acciones y búsqueda">
          {view !== "audit" && <label className="platform-search"><Search aria-hidden="true" /><span className="sr-only">Buscar</span><input value={query} onChange={(event) => setQuery(event.target.value)} placeholder={view === "users" ? "Buscar usuario, correo u organización" : "Buscar organización o responsable"} /></label>}
          <div className="toolbar-actions">
            <button className="secondary-action" type="button" onClick={() => setPanel(panel === "invitation" ? null : "invitation")}><UserPlus aria-hidden="true" />Invitar usuario</button>
            <button className="primary-action" type="button" onClick={() => setPanel(panel === "tenant" ? null : "tenant")}><Plus aria-hidden="true" />Nueva organización</button>
          </div>
        </section>

        {panel && <section className="platform-action-panel" aria-label={panel === "tenant" ? "Crear organización" : "Invitar usuario"}>
          {panel === "tenant" ? <TenantForm busy={busy === "create-tenant"} onSubmit={handleCreateTenant} onClose={() => setPanel(null)} /> : <InvitationForm tenants={snapshot.tenants} busy={busy === "invite-user"} onSubmit={handleInvite} onClose={() => setPanel(null)} />}
        </section>}

        {view === "overview" && <div className="overview-grid"><TenantTable tenants={filteredTenants.slice(0, 8)} busy={busy} onStatus={(tenant, status) => runMutation(`tenant-${tenant.id}`, () => setTenantStatus(tenant.id, status), `Estado de ${tenant.name} actualizado.`)} /><AuditRail events={snapshot.audit.slice(0, 7)} /></div>}
        {view === "tenants" && <TenantTable tenants={filteredTenants} busy={busy} onStatus={(tenant, status) => runMutation(`tenant-${tenant.id}`, () => setTenantStatus(tenant.id, status), `Estado de ${tenant.name} actualizado.`)} />}
        {view === "users" && <UserTable users={filteredUsers} currentUserId={user.id} busy={busy} onStatus={(entry, status) => runMutation(`user-${entry.id}`, () => setUserStatus(entry.id, status), `Estado de ${entry.displayName} actualizado.`)} onRevoke={(entry) => runMutation(`sessions-${entry.id}`, () => revokeUserSessions(entry.id), `Sesiones de ${entry.displayName} revocadas.`)} onRole={(entry, membership, role) => runMutation(`role-${entry.id}-${membership.tenantId}`, () => setMembershipRole(membership.tenantId, entry.id, role), `Rol de ${entry.displayName} actualizado.`)} />}
        {view === "audit" && <AuditTable events={snapshot.audit} />}
      </main>
    </div>
  );
}

function NavButton({ active, onClick, icon, label }: { active: boolean; onClick: () => void; icon: React.ReactNode; label: string }) {
  return <button className={active ? "is-active" : ""} type="button" onClick={onClick}>{icon}<span>{label}</span></button>;
}

function TenantForm({ busy, onSubmit, onClose }: { busy: boolean; onSubmit: (event: FormEvent<HTMLFormElement>) => void; onClose: () => void }) {
  return <form onSubmit={onSubmit}><div className="panel-heading"><div><h2>Nueva organización</h2><p>Crea el espacio del cliente; después podrás asignar su propietario.</p></div><button type="button" onClick={onClose}>Cerrar</button></div><div className="form-row"><label>Nombre de la organización<input name="name" required maxLength={160} placeholder="Ej. Fundación Illimani" /></label><label>Identificador<input name="slug" required maxLength={80} pattern="[a-z0-9]+(?:-[a-z0-9]+)*" placeholder="fundacion-illimani" /><small>Minúsculas, números y guiones.</small></label><button className="primary-action" disabled={busy}>{busy ? "Creando…" : "Crear organización"}</button></div></form>;
}

function InvitationForm({ tenants, busy, onSubmit, onClose }: { tenants: TenantRecord[]; busy: boolean; onSubmit: (event: FormEvent<HTMLFormElement>) => void; onClose: () => void }) {
  return <form onSubmit={onSubmit}><div className="panel-heading"><div><h2>Invitar usuario</h2><p>Asigna la organización y el rol desde el inicio.</p></div><button type="button" onClick={onClose}>Cerrar</button></div>{tenants.length === 0 ? <div className="inline-empty">Primero crea una organización para poder asignar usuarios.</div> : <div className="form-row invite-row"><label>Organización<select name="tenantId" required defaultValue=""><option value="" disabled>Selecciona una organización</option>{tenants.filter((tenant) => tenant.status !== "deactivated").map((tenant) => <option key={tenant.id} value={tenant.id}>{tenant.name}</option>)}</select></label><label>Nombre completo<input name="displayName" required maxLength={160} placeholder="Nombre de la persona" /></label><label>Correo electrónico<input name="email" required type="email" placeholder="nombre@organizacion.bo" /></label><label>Rol<select name="role" defaultValue="tenant_owner"><option value="tenant_owner">Propietario</option><option value="tenant_admin">Administrador</option><option value="guardian">Tutor</option></select></label><button className="primary-action" disabled={busy}>{busy ? "Preparando…" : "Crear invitación"}</button></div>}</form>;
}

function TenantTable({ tenants, busy, onStatus }: { tenants: TenantRecord[]; busy: string | null; onStatus: (tenant: TenantRecord, status: TenantRecord["status"]) => void }) {
  if (tenants.length === 0) return <EmptyState icon={<Building2 />} title="Aún no hay organizaciones" text="Crea la primera organización y luego invita a su propietario." />;
  return <section className="data-region"><div className="region-heading"><h2>Organizaciones</h2><span>{tenants.length} registros</span></div><div className="table-scroll"><table className="tenant-table"><thead><tr><th>Organización</th><th>Responsable</th><th>Usuarios</th><th>Sesiones</th><th>Estado y acción</th></tr></thead><tbody>{tenants.map((tenant) => <tr key={tenant.id}><td><strong>{tenant.name}</strong><small>{tenant.slug}</small></td><td>{tenant.ownerName ? <><span>{tenant.ownerName}</span><small>{tenant.ownerEmail}</small></> : <span className="muted">Sin propietario</span>}</td><td className="tabular">{tenant.userCount}</td><td className="tabular">{tenant.activeSessions}</td><td><form className="row-action" onSubmit={(event) => { event.preventDefault(); const status = String(new FormData(event.currentTarget).get("status")) as TenantRecord["status"]; onStatus(tenant, status); }}><span className={`status-mark is-${tenant.status}`}>{statusLabels[tenant.status]}</span><select name="status" defaultValue={tenant.status} aria-label={`Estado de ${tenant.name}`}><option value="active">Activo</option><option value="suspended">Suspendido</option><option value="deactivated">Desactivado</option></select><button disabled={busy === `tenant-${tenant.id}`}>{busy === `tenant-${tenant.id}` ? "Guardando…" : "Guardar"}</button></form></td></tr>)}</tbody></table></div></section>;
}

function UserTable({ users, currentUserId, busy, onStatus, onRevoke, onRole }: { users: PlatformSnapshot["users"]; currentUserId: string; busy: string | null; onStatus: (user: PlatformSnapshot["users"][number], status: "active" | "suspended" | "deactivated") => void; onRevoke: (user: PlatformSnapshot["users"][number]) => void; onRole: (user: PlatformSnapshot["users"][number], membership: Membership, role: Membership["role"]) => void }) {
  if (users.length === 0) return <EmptyState icon={<Users />} title="Aún no hay usuarios" text="Invita al propietario de una organización para comenzar." />;
  return <section className="data-region"><div className="region-heading"><h2>Usuarios y permisos</h2><span>{users.length} registros</span></div><div className="table-scroll"><table className="user-table"><thead><tr><th>Usuario</th><th>Organización y rol</th><th>Sesiones</th><th>Estado</th><th>Control</th></tr></thead><tbody>{users.map((entry) => <tr key={entry.id}><td><strong>{entry.displayName}</strong><small>{entry.email}</small>{entry.platformRole && <small className="platform-role">{entry.platformRole}</small>}</td><td>{entry.memberships.length === 0 ? <span className="muted">Sin organización</span> : entry.memberships.map((membership) => <form className="membership-line" key={membership.tenantId} onSubmit={(event) => { event.preventDefault(); const role = String(new FormData(event.currentTarget).get("role")) as Membership["role"]; onRole(entry, membership, role); }}><span>{membership.tenantName}</span><select name="role" defaultValue={membership.role} aria-label={`Rol de ${entry.displayName} en ${membership.tenantName}`}><option value="tenant_owner">Propietario</option><option value="tenant_admin">Administrador</option><option value="guardian">Tutor</option></select><button disabled={busy === `role-${entry.id}-${membership.tenantId}`}>Guardar rol</button></form>)}</td><td className="tabular">{entry.activeSessions}</td><td><span className={`status-mark is-${entry.status}`}>{statusLabels[entry.status]}</span></td><td><form className="user-controls" onSubmit={(event) => { event.preventDefault(); const status = String(new FormData(event.currentTarget).get("status")) as "active" | "suspended" | "deactivated"; onStatus(entry, status); }}><select name="status" defaultValue={entry.status === "invited" ? "suspended" : entry.status} disabled={entry.id === currentUserId || entry.status === "invited"} aria-label={`Cambiar estado de ${entry.displayName}`}><option value="active">Activo</option><option value="suspended">Suspendido</option><option value="deactivated">Desactivado</option></select><button disabled={entry.id === currentUserId || entry.status === "invited" || busy === `user-${entry.id}`}>Guardar estado</button><button className="danger-text" type="button" disabled={entry.activeSessions === 0 || busy === `sessions-${entry.id}`} onClick={() => onRevoke(entry)}><KeyRound aria-hidden="true" />Revocar sesiones</button></form></td></tr>)}</tbody></table></div></section>;
}

function AuditRail({ events }: { events: AuditEvent[] }) {
  return <aside className="audit-rail"><div className="region-heading"><h2>Actividad reciente</h2><ShieldCheck aria-hidden="true" /></div>{events.length === 0 ? <p className="muted">Las acciones administrativas aparecerán aquí.</p> : events.map((event) => <AuditItem key={event.id} event={event} />)}</aside>;
}

function AuditTable({ events }: { events: AuditEvent[] }) {
  if (events.length === 0) return <EmptyState icon={<Activity />} title="Sin registros de auditoría" text="Las acciones de plataforma se registrarán automáticamente." />;
  return <section className="data-region"><div className="region-heading"><h2>Registro de auditoría</h2><span>{events.length} eventos recientes</span></div><div className="audit-list">{events.map((event) => <AuditItem key={event.id} event={event} detailed />)}</div></section>;
}

function AuditItem({ event, detailed = false }: { event: AuditEvent; detailed?: boolean }) {
  const when = new Intl.DateTimeFormat("es-BO", { dateStyle: detailed ? "medium" : undefined, timeStyle: "short", timeZone: "America/La_Paz" }).format(new Date(event.occurredAt));
  return <article className="audit-event"><span className="registration-point" aria-hidden="true" /><div><strong>{auditLabels[event.action] ?? event.action}</strong><p>{event.tenantName ?? event.actorEmail ?? "Sistema APU"}</p>{detailed && <small>{event.actorEmail ?? "Proceso del sistema"} · {event.resourceType}</small>}</div><time dateTime={event.occurredAt}>{when}</time></article>;
}

function EmptyState({ icon, title, text }: { icon: React.ReactNode; title: string; text: string }) {
  return <section className="empty-state">{icon}<h2>{title}</h2><p>{text}</p></section>;
}

function viewTitle(view: View) {
  return { overview: "Operación de plataforma", tenants: "Organizaciones", users: "Usuarios y permisos", audit: "Auditoría operacional" }[view];
}
