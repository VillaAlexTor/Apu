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
      </section>

      <aside className="registration-stamp" aria-label="Origen de APU">
        <span>Registro cartográfico</span>
        <small>APU / Acceso 001</small>
      </aside>

      <section className="login-legend" aria-label="Formulario de acceso">
        <span className="route-arrival" aria-hidden="true"><span /></span>
        <LoginForm />
      </section>
    </main>
  );
}
