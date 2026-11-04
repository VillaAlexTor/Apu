# Arquitectura SaaS multiusuario

**Estado:** propuesta de producto privado  
**Última actualización:** 2026-11-04

## 1. Modelo de servicio

APU será un servicio SaaS privado operado por el propietario del producto. El
operador crea y administra cuentas de clientes; cada cliente paga una
suscripción mensual para acceder al panel y administrar sus dispositivos.

```text
Operador APU
└── Tenant / Cliente
    ├── Suscripción
    ├── Usuarios y roles
    ├── Familias
    │   └── Dispositivos Android
    └── Datos, reglas y auditoría
```

`tenant_id` es el límite obligatorio de aislamiento. `family_id` y `device_id`
son límites subordinados y nunca sustituyen la comprobación del tenant.

## 2. Componentes

- **Panel de clientes:** gestión de usuarios autorizados, familias, módulos y
  dispositivos.
- **Panel administrativo:** alta, suspensión y soporte operacional de tenants,
  sin acceso ordinario a coordenadas descifradas.
- **API:** autenticación, autorización, facturación, dispositivos, datos y
  auditoría.
- **Worker:** retención, renovaciones, notificaciones y tareas asíncronas.
- **Proveedor de pagos:** fuente externa de estados de pago; la API conserva
  referencias, no datos completos de tarjetas.
- **Android:** agente visible y consentido del dispositivo.
- **PostgreSQL:** fuente de verdad con partición lógica por tenant.

## 3. Regla de aislamiento

Toda consulta y mutación debe seguir:

```text
actor → tenant → family → device
```

El backend obtiene el tenant desde la sesión y nunca confía en un `tenant_id`
enviado por el navegador. Los repositorios reciben el contexto autorizado o
aplican filtros obligatorios. Los tests deben demostrar que un ID válido de otro
tenant no devuelve datos ni modifica recursos.

## 4. Suscripción como capacidad

La suscripción habilita el acceso comercial del tenant. Su estado se evalúa en
el backend para cada operación que consume el servicio:

- `trial`: acceso temporal según política;
- `active`: acceso normal;
- `past_due`: acceso limitado durante periodo de gracia;
- `suspended`: no permite nuevas operaciones;
- `cancelled`: termina al cierre del periodo;
- `terminated`: acceso cerrado y datos sujetos a retención.

La suspensión no debe borrar inmediatamente datos ni desvincular dispositivos.
Debe existir un proceso de recuperación y exportación según los términos
aplicables.

## 5. Límites operativos

Cada plan define cantidad de usuarios, dispositivos, retención, frecuencia de
sincronización y módulos habilitados. Los límites se validan en el caso de uso,
no solo en la interfaz.

## 6. Seguridad del operador

El personal interno usa cuentas separadas, MFA, roles mínimos y auditoría.
Soporte puede consultar estados técnicos y metadatos mínimos, pero no leer
coordenadas o contenido sensible salvo un procedimiento excepcional autorizado y
registrado.

