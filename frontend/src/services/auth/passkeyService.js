/**
 * passkeyService.js
 * Maneja toda la comunicación con el backend para WebAuthn/Passkey,
 * incluyendo la codificación/decodificación de ArrayBuffer ↔ base64url
 * que requiere la Web Authentication API.
 */
import { getApiBaseUrl } from '@/config/api.js'

const BASE = getApiBaseUrl()

// ── Utilidades de codificación ────────────────────────────────────────────

/**
 * Convierte un ArrayBuffer o Uint8Array a una cadena base64url
 * (sin padding, compatible con WebAuthn)
 */
function bufferToBase64url(buffer) {
  const bytes = buffer instanceof Uint8Array ? buffer : new Uint8Array(buffer)
  let str = ''
  for (const b of bytes) str += String.fromCharCode(b)
  return btoa(str).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
}

/**
 * Decodifica una cadena base64url a Uint8Array
 */
function base64urlToBuffer(b64url) {
  const b64 = b64url.replace(/-/g, '+').replace(/_/g, '/').padEnd(
    b64url.length + ((4 - (b64url.length % 4)) % 4),
    '='
  )
  const bin = atob(b64)
  const buf = new Uint8Array(bin.length)
  for (let i = 0; i < bin.length; i++) buf[i] = bin.charCodeAt(i)
  return buf
}

/**
 * Convierte los campos ArrayBuffer de las opciones de creación (registro)
 * al formato que espera navigator.credentials.create()
 */
function prepareRegistrationOptions(options) {
  const pubKey = options.publicKey
  return {
    publicKey: {
      ...pubKey,
      challenge: base64urlToBuffer(pubKey.challenge),
      user: {
        ...pubKey.user,
        id: base64urlToBuffer(pubKey.user.id),
      },
      excludeCredentials: (pubKey.excludeCredentials || []).map(c => ({
        ...c,
        id: base64urlToBuffer(c.id),
      })),
    },
  }
}

/**
 * Convierte los campos ArrayBuffer de las opciones de aserción (login)
 * al formato que espera navigator.credentials.get()
 */
function prepareLoginOptions(options) {
  const pubKey = options.publicKey
  return {
    publicKey: {
      ...pubKey,
      challenge: base64urlToBuffer(pubKey.challenge),
      allowCredentials: (pubKey.allowCredentials || []).map(c => ({
        ...c,
        id: base64urlToBuffer(c.id),
      })),
    },
  }
}

/**
 * Serializa la respuesta de navigator.credentials.create() para enviarla al backend
 */
function serializeRegistrationCredential(credential) {
  const resp = credential.response
  return {
    id: credential.id,
    rawId: bufferToBase64url(credential.rawId),
    type: credential.type,
    response: {
      clientDataJSON: bufferToBase64url(resp.clientDataJSON),
      attestationObject: bufferToBase64url(resp.attestationObject),
      transports: resp.getTransports ? resp.getTransports() : [],
      publicKeyAlgorithm: resp.getPublicKeyAlgorithm ? resp.getPublicKeyAlgorithm() : -7,
      publicKey: resp.getPublicKey ? bufferToBase64url(resp.getPublicKey()) : null,
      authenticatorData: resp.getAuthenticatorData ? bufferToBase64url(resp.getAuthenticatorData()) : null,
    },
    authenticatorAttachment: credential.authenticatorAttachment || null,
    clientExtensionResults: credential.getClientExtensionResults
      ? credential.getClientExtensionResults()
      : {},
  }
}

/**
 * Serializa la respuesta de navigator.credentials.get() para enviarla al backend
 */
function serializeLoginCredential(credential) {
  const resp = credential.response
  return {
    id: credential.id,
    rawId: bufferToBase64url(credential.rawId),
    type: credential.type,
    response: {
      clientDataJSON: bufferToBase64url(resp.clientDataJSON),
      authenticatorData: bufferToBase64url(resp.authenticatorData),
      signature: bufferToBase64url(resp.signature),
      userHandle: resp.userHandle ? bufferToBase64url(resp.userHandle) : null,
    },
    authenticatorAttachment: credential.authenticatorAttachment || null,
    clientExtensionResults: credential.getClientExtensionResults
      ? credential.getClientExtensionResults()
      : {},
  }
}

// ── API calls ─────────────────────────────────────────────────────────────

async function apiFetch(path, { method = 'POST', body = null, token = null, query = '' } = {}) {
  const url = `${BASE}${path}${query ? '?' + query : ''}`
  const headers = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(url, {
    method,
    headers,
    body: body !== null ? JSON.stringify(body) : undefined,
  })
  const data = await res.json()
  if (!res.ok || !data.success) {
    throw new Error(data.error || `Error ${res.status}`)
  }
  return data.data
}

// ── Registro ──────────────────────────────────────────────────────────────

/**
 * Registra una nueva passkey para el usuario autenticado.
 * @param {string} token  JWT del usuario
 * @param {string} name   Nombre descriptivo de la passkey
 */
export async function registerPasskey(token, name = 'Mi Passkey') {
  if (!window.PublicKeyCredential) {
    throw new Error('Tu navegador no soporta Passkeys')
  }

  // 1. Obtener opciones del servidor
  const options = await apiFetch('/auth/passkey/register/begin', { token })

  // 2. Invocar el autenticador del dispositivo
  const credential = await navigator.credentials.create(prepareRegistrationOptions(options))
  if (!credential) throw new Error('Registro cancelado por el usuario')

  // 3. Enviar respuesta al servidor
  const serialized = serializeRegistrationCredential(credential)
  await apiFetch('/auth/passkey/register/finish', {
    token,
    body: serialized,
    query: `name=${encodeURIComponent(name)}`,
  })
}

// ── Login ─────────────────────────────────────────────────────────────────

/**
 * Intenta autenticar al usuario con el método biométrico/PIN del dispositivo.
 * Usa mediation:'optional' para que el navegador dispare la UI nativa del dispositivo
 * si hay credenciales disponibles, o resuelva null silenciosamente si no hay ninguna.
 * @returns {{ token, user, expires_in, token_type }}
 */
export async function loginWithPasskey() {
  if (!window.PublicKeyCredential) {
    throw new Error('NO_SUPPORT')
  }

  // 1. Obtener el challenge del servidor
  const { options, session_id } = await apiFetch('/auth/passkey/login/begin')

  // 2. Solicitar autenticación al dispositivo con mediation:'optional'
  //    → Si hay credenciales: muestra la UI nativa (biometría/PIN) automáticamente
  //    → Si no hay credenciales: resuelve a null sin mostrar nada
  const credential = await navigator.credentials.get({
    ...prepareLoginOptions(options),
    mediation: 'optional',
  })

  if (!credential) {
    throw new Error('NO_CREDENTIALS')
  }

  // 3. Verificar en el servidor y obtener JWT
  const serialized = serializeLoginCredential(credential)
  return apiFetch('/auth/passkey/login/finish', {
    body: serialized,
    query: `session_id=${session_id}`,
  })
}

// ── Gestión ───────────────────────────────────────────────────────────────

/**
 * Lista las passkeys del usuario autenticado.
 * @param {string} token
 */
export async function listPasskeys(token) {
  return apiFetch('/auth/passkey/credentials', { method: 'GET', token })
}

/**
 * Elimina una passkey del usuario autenticado.
 * @param {string} token
 * @param {number} id
 */
export async function deletePasskey(token, id) {
  return apiFetch(`/auth/passkey/credentials/${id}`, { method: 'DELETE', token })
}
