import Image from "next/image";
import { LoginForm } from "@/components/auth/login-form";
import { CartographicFrame } from "@/components/cartographic/cartographic-frame";

export default function Home() {
  return (
    <main className="login-page">
      <Image className="map-background" src="/assets/apu-topographic-field.webp" alt="" fill priority sizes="100vw" />
      <div className="map-quiet-field" aria-hidden="true" />
      <CartographicFrame />

      <section className="brand-block" aria-labelledby="login-title">
        <div className="brand-lockup">
          <span className="wordmark">APU</span>
          <span className="survey-rule" aria-hidden="true" />
        </div>
        <h1 id="login-title">Cuidar también es ser transparentes.</h1>
        <p>Un solo acceso para cada persona autorizada. Tus permisos definen el espacio que verás después.</p>
      </section>

      <aside className="registration-stamp" aria-label="Origen de APU">
        <span>Registro cartográfico</span>
        <strong>La Paz · Bolivia</strong>
        <small>APU / Acceso 001</small>
      </aside>

      <section className="login-legend" aria-label="Formulario de acceso">
        <span className="route-arrival" aria-hidden="true"><span /></span>
        <LoginForm />
      </section>
    </main>
  );
}
